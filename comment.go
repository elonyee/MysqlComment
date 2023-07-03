package mysql_comment

import (
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"regexp"
	"strings"
)

type DbConfig struct {
	DSN      string
	Host     string
	Port     int
	Database string
	Username string
	Password string
	Prefix   string
	Charset  string
}

type ResetComment struct {
	DB     *gorm.DB
	Config DbConfig
	oldSql []string
	newSql []string
}

func (r *ResetComment) Connect() *ResetComment {
	if r.Config.DSN == "" {
		r.Config.DSN = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			r.Config.Username,
			r.Config.Password,
			r.Config.Host,
			r.Config.Port,
			r.Config.Database,
			r.Config.Charset,
		)
	}
	strategySet := schema.NamingStrategy{
		TablePrefix:   r.Config.Prefix,
		SingularTable: true,
	}
	var loggerSet logger.Interface
	db, err := gorm.Open(mysql.Open(r.Config.DSN), &gorm.Config{
		NamingStrategy:                           strategySet,
		Logger:                                   loggerSet,
		DisableForeignKeyConstraintWhenMigrating: true,
		SkipDefaultTransaction:                   false,
	})
	if err != nil {
		panic(err)
	}
	r.DB = db
	return r
}

// CreateSql 生成SQL 自己手动执行（请自行检查生成的SQL是否符合预期）
func (r *ResetComment) CreateSql() string {
	r.newSql = r.getList()
	return r.ToSqlString()
}

// RemoveComment 直接执行清空Comment操作（注意：此操作风险非常大！可能会损坏数据表结构！请在测试环境测试没问题后再进行！）
func (r *ResetComment) RemoveComment() (result string, err error) {
	var runInfo []map[string]interface{}
	r.newSql = r.getList()
	for _, v := range r.newSql {
		err = r.DB.Raw(v).Scan(&runInfo).Error
		if err != nil {
			err = errors.New(fmt.Sprintf("Error executing SQL [ %s ]:%s \r\n", v, err.Error()))
			break
		}
	}
	if err != nil {
		result = "discontinue."
	} else {
		result = "succeed."
	}
	return
}

func (r *ResetComment) getList() []string {
	var tableList []map[string]interface{}
	var sql []string
	err := r.DB.Raw("show tables").Scan(&tableList).Error
	if err != nil {
		return sql
	}
	preLen := len(r.Config.Prefix)
	for _, v := range tableList {
		for _, vv := range v {
			tableName := fmt.Sprintf("%s", vv)
			if preLen > 0 {
				if len(tableName) <= preLen || r.Config.Prefix != tableName[:preLen] {
					continue
				}
			}
			info := r.getInfo(tableName)
			if len(info) > 0 {
				sql = append(sql, info...)
			}
		}
	}
	return sql
}

func (r *ResetComment) getInfo(tableName string) []string {
	var tableInfo map[string]interface{}
	var sql []string
	err := r.DB.Raw(fmt.Sprintf("SHOW CREATE TABLE `%s`", tableName)).Scan(&tableInfo).Error
	if err != nil {
		return sql
	}
	fields := strings.Split(fmt.Sprintf("%s", tableInfo["Create Table"]), "\n")
	count := len(fields)
	if count > 1 {
		if strings.Index(fmt.Sprintf("%s", tableInfo["Create Table"]), " COMMENT=") > 0 {
			sql = append(sql, fmt.Sprintf("\nALTER TABLE `%s` COMMENT '';", tableInfo["Table"]))
			r.oldSql = append(r.oldSql, fmt.Sprintf("\nALTER TABLE `%s` COMMENT '%s';", tableInfo["Table"], r.getTableComment(fmt.Sprintf("%s;", tableInfo["Create Table"]))))
		} else {
			sql = append(sql, fmt.Sprintf("\nALTER TABLE `%s` COMMENT '';", tableInfo["Table"]))
			r.oldSql = append(r.oldSql, fmt.Sprintf("\nALTER TABLE `%s` COMMENT '';", tableInfo["Table"]))
		}
	}
	for _, v := range fields {
		line := strings.Split(strings.TrimSpace(v), " COMMENT '")
		if line[0][:1] != "`" || len(line) > 2 {
			continue
		}
		if len(line) < 2 {
			line = append(line, "")
			if line[0][len(line[0])-1:] == "," {
				line[0] = line[0][0 : len(line[0])-1]
			}
		}
		sql = append(sql, fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s COMMENT '';", tableInfo["Table"], line[0]))
		r.oldSql = append(r.oldSql, fmt.Sprintf("ALTER TABLE `%s` MODIFY COLUMN %s COMMENT '%s';", tableInfo["Table"], line[0], r.getFieldComment(fmt.Sprintf("%s,", v))))
	}
	return sql
}

func (r *ResetComment) getTableComment(sql string) string {
	sql = strings.ReplaceAll(sql, "\r", "")
	sql = strings.ReplaceAll(sql, "\n", " ")
	exp := regexp.MustCompile(`( COMMENT='(.*)';)+$`)
	str := exp.FindStringSubmatch(sql)
	if len(str) == 3 {
		return strings.ReplaceAll(str[2], "'", " ")
	}
	return ""
}

func (r *ResetComment) getFieldComment(sql string) string {
	sql = strings.ReplaceAll(sql, "\r", "")
	sql = strings.ReplaceAll(sql, "\n", " ")
	exp := regexp.MustCompile(` COMMENT '(.*)',+$`)
	str := exp.FindStringSubmatch(sql)
	if len(str) == 2 {
		return strings.ReplaceAll(str[1], "'", " ")
	}
	return ""
}

func (r *ResetComment) ToSqlString() string {
	return fmt.Sprintf("# -- -------------- Comment Restore SQL Start ---------------------\n%s \n# -- -------------- Comment Restore SQL End ---------------------\n\n# -- -------------- Comment Remove SQL Start ---------------------\n%s\n# -- -------------- Comment Remove SQL End ---------------------\n", strings.Join(r.oldSql, "\n"), strings.Join(r.newSql, "\n"))
}
