package main

import (
	"github.com/hezhis/go_autodb"
	"github.com/hezhis/go_autodb/db"
	logger "github.com/hezhis/go_log"
)

func main() {
	if err := db.InitOrmMysql("game", "game@2021", "127.0.0.1", 3306, "cq_charge"); nil != err {
		logger.Fatal("%v", err)
	}

	if err := go_autodb.ParseTableConf(); nil != err {
		logger.Fatal("%v", err)
	}
	if err := go_autodb.BuildTables(); nil != err {
		logger.Fatal("%v", err)
	}
}
