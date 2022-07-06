package go_autodb

import (
	"log"

	"github.com/hezhis/go_autodb/db"
)

func BuildTables() error {
	for _, table := range tables {
		if err := table.Build(); nil != err {
			return err
		}
	}
	for _, procedure := range Procedures {
		if err := procedure.Build(); nil != err {
			return err
		}
	}

	//execSQL("call initdb", false)
	return nil
}

func execSQL(sql string, echo bool) error {
	if len(sql) <= 0 {
		return nil
	}
	if echo {
		log.Printf("%s\n", sql)
	}
	_, err := db.OrmEngine.Exec(sql)
	return err
}
