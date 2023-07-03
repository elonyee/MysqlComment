# Mysql Comment

生成 Mysql 表结构注释 SQL。批量增加、批量移除 Mysql 表结构注释。

Mysql comment batch increase and batch removal.



## 生成 Mysql 表结构注释 SQL 使用示例`./tests/CreateSqlTest.go`：
```golang
package main

import (
	"fmt"
	"github.com/zaesn/mysql_comment"
)

func main() {
	conf := mysql_comment.DbConfig{
		Host:     "127.0.0.1",
		Port:     3306,
		Username: "root",
		Password: "password",
		Database: "test_database",
		Prefix:   "",
		Charset:  "utf8mb4",
	}
	resetComment := &mysql_comment.ResetComment{Config: conf}
	result := resetComment.Connect().CreateSql()
	fmt.Println(result)
}
```
- 输出:
```sh
> go run ./CreateSqlTest.go
# -- -------------- Comment Restore SQL Start ---------------------

ALTER TABLE `district` COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '区划信息id';
ALTER TABLE `district` MODIFY COLUMN `pid` int(11) DEFAULT NULL COMMENT '父级挂接id';
ALTER TABLE `district` MODIFY COLUMN `code` varchar(255) DEFAULT NULL COMMENT '区划编码';
ALTER TABLE `district` MODIFY COLUMN `name` varchar(255) DEFAULT NULL COMMENT '区划名称';
ALTER TABLE `district` MODIFY COLUMN `remark` varchar(255) DEFAULT NULL COMMENT '备注';
ALTER TABLE `district` MODIFY COLUMN `create_time` datetime DEFAULT NULL COMMENT '创建时间';
ALTER TABLE `district` MODIFY COLUMN `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';
ALTER TABLE `district` MODIFY COLUMN `status` tinyint(1) DEFAULT NULL COMMENT '状态 0 正常 -2 删除 -1 停用';
ALTER TABLE `district` MODIFY COLUMN `level` tinyint(1) DEFAULT NULL COMMENT '级次id 0:省/自治区/直辖市 1:市级 2:县级'; 
# -- -------------- Comment Restore SQL End ---------------------

# -- -------------- Comment Remove SQL Start ---------------------

ALTER TABLE `district` COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `pid` int(11) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `code` varchar(255) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `name` varchar(255) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `remark` varchar(255) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `create_time` datetime DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `status` tinyint(1) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `level` tinyint(1) DEFAULT NULL COMMENT '';
# -- -------------- Comment Remove SQL End ---------------------

```



## 直接执行清空Comment操作 使用示例`./tests/RemoveCommentTest.go`：
- 注意：此操作风险非常大！可能会损坏数据表结构！请在测试环境测试没问题后再进行！

```golang
package main

import (
	"fmt"
	"github.com/zaesn/mysql_comment"
)

func main() {
	conf := mysql_comment.DbConfig{DSN: fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		"root",
		"password",
		"localhost",
		3306,
		"test_database",
		"utf8mb4",
	)}
	resetComment := &mysql_comment.ResetComment{Config: conf}
	result, err := resetComment.Connect().RemoveComment()
	fmt.Printf("sql: \n %s \n", resetComment.ToSqlString())
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(result)
}
```
- 输出:
```sh
> go run ./RemoveCommentTest.go
sql: 
 # -- -------------- Comment Restore SQL Start ---------------------

ALTER TABLE `district` COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '区划信息id';
ALTER TABLE `district` MODIFY COLUMN `pid` int(11) DEFAULT NULL COMMENT '父级挂接id';
ALTER TABLE `district` MODIFY COLUMN `code` varchar(255) DEFAULT NULL COMMENT '区划编码';
ALTER TABLE `district` MODIFY COLUMN `name` varchar(255) DEFAULT NULL COMMENT '区划名称';
ALTER TABLE `district` MODIFY COLUMN `remark` varchar(255) DEFAULT NULL COMMENT '备注';
ALTER TABLE `district` MODIFY COLUMN `create_time` datetime DEFAULT NULL COMMENT '创建时间';
ALTER TABLE `district` MODIFY COLUMN `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间';
ALTER TABLE `district` MODIFY COLUMN `status` tinyint(1) DEFAULT NULL COMMENT '状态 0 正常 -2 删除 -1 停用';
ALTER TABLE `district` MODIFY COLUMN `level` tinyint(1) DEFAULT NULL COMMENT '级次id 0:省/自治区/直辖市 1:市级 2:县级';
# -- -------------- Comment Restore SQL End ---------------------

# -- -------------- Comment Remove SQL Start ---------------------

ALTER TABLE `district` COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `pid` int(11) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `code` varchar(255) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `name` varchar(255) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `remark` varchar(255) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `create_time` datetime DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `update_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `status` tinyint(1) DEFAULT NULL COMMENT '';
ALTER TABLE `district` MODIFY COLUMN `level` tinyint(1) DEFAULT NULL COMMENT '';
# -- -------------- Comment Remove SQL End ---------------------

succeed.

```


- 用于测试的`district`数据表结构来源：`https://github.com/MingA21/District-SQL`