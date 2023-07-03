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
