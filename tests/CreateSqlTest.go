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
