package go_autodb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/hezhis/go_autodb/column"
)

var (
	Tables     []*TableConf
	Procedures []*ProcedureSt
)

type (
	ColumnConf struct {
		Name          string
		Type          string
		Comment       string
		AutoIncrement bool
		Size          int
		Unsigned      bool
		Default       string
	}

	TableConf struct {
		Name    string
		Comment string
		Columns []*ColumnConf
		Keys    []*TableKey
	}
)

func loadTableConf(file string) bool {
	data, err := ioutil.ReadFile(file)
	if nil != err {
		log.Printf("load table conf error! [%s] %s\n", file, err)
		return false
	}

	if err = json.Unmarshal(data, Tables); err != nil {
		log.Fatalf("load %s Unmarshal json error:%s", file, err)
		return false
	}
	return true
}

func ParseTableConf() error {
	for _, line := range Tables {
		table := NewTable(line.Name, line.Comment, len(line.Columns))
		for _, field := range line.Columns {
			var col column.IColumn
			var err error
			if field.Type == "varchar" {
				if col, err = column.NewStringColumn(field.Name, field.Size, field.Comment, field.Default); nil != err {
					return err
				}
			} else if field.Type == "datetime" {
				if col, err = column.NewDateTimeColumn(field.Name, field.Comment); nil != err {
					return err
				}
			} else if strings.Contains(field.Type, "blob") {
				if col, err = column.NewBlobColumn(field.Name, field.Type, field.Comment); nil != err {
					return err
				}
			} else if strings.Contains(field.Type, "int") {
				if col, err = column.NewIntColumn(field.Name, field.Type, field.Unsigned, field.Comment, field.AutoIncrement, field.Default); nil != err {
					return err
				}
			} else {
				return errors.New(fmt.Sprintf("the type of column %s of table %s ", line.Name, field.Name))
			}
			table.addColumn(col)
		}
		for _, key := range line.Keys {
			table.tblKeys[key.Name] = key.Copy()
		}
		tables[line.Name] = table
	}
	return nil
}
