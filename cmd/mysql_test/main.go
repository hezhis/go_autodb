package main

import (
	"log"

	"github.com/hezhis/go_autodb"
	"github.com/hezhis/go_autodb/db"
)

func main() {
	if err := db.InitOrmMysql("server", "123456abc", "127.0.0.1", 3306, "sh3d_log"); nil != err {
		log.Fatalln(err)
	}

	if err := go_autodb.ParseTableConf(); nil != err {
		log.Fatalln(err)
	}
	if err := go_autodb.BuildTables(); nil != err {
		log.Fatalln(err)
	}
}
