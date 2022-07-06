package go_autodb

import (
	"errors"
	"fmt"

	"github.com/hezhis/go_autodb/column"
)

func (st *Table) checkConf() error {
	for name, key := range st.tblKeys {
		for _, column := range key.ColumnVec {
			if _, exists := st.Columns[column]; !exists {
				return errors.New(fmt.Sprintf("index error! column of table[%s] is not exists! key[%s]", st.Name, name))
			}
		}
	}
	return nil
}

func (st *Table) checkCompatible() error {
	sqlInfo, err := st.GetTableColumnInfo(false)
	if nil != err {
		return err
	}
	if nil == sqlInfo {
		return errors.New(fmt.Sprintf("%s table column info is nil", st.Name))
	}

	// check column compatible
	for _, info := range sqlInfo.Columns {
		if column, exists := st.Columns[info.Field]; exists {
			if column.IsEqual(info) {
				continue
			}
			if !column.IsCompatible(info) {
				return errors.New(fmt.Sprintf("column[%s] of table [%s] is uncompatible", column.GetName(), st.Name))
			}
		} else {
			// no empty table forbid delete column
			// only check here, delete in change function
			flag, err := st.hasData()
			if nil != err {
				return err
			}
			if flag {
				return errors.New(fmt.Sprintf("table[%s] is no empty, delete column is forbiden!", st.Name))
			}
		}
	}
	// check index compatible
	for _, info := range sqlInfo.Keys {
		newKey, ok := st.tblKeys[info.Name]
		if !ok {
			return errors.New(fmt.Sprintf("delete exists key is forbiden. table[%s], key[%s]", st.Name, info.Name))
		}
		if info.Type != newKey.Type {
			return errors.New(fmt.Sprintf("change exists key is forbiden. table[%s], key[%s]", st.Name, info.Name))
		}
		size := len(info.ColumnVec)
		if len(newKey.ColumnVec) != size {
			return errors.New(fmt.Sprintf("the config key is unmatch with db. table[%s], key[%s]", st.Name, info.Name))
		}
		for i := 0; i < size; i++ {
			if info.ColumnVec[i] != newKey.ColumnVec[i] {
				return errors.New(fmt.Sprintf("change exists key is forbiden. table[%s], key[%s]", st.Name, info.Name))
			}
		}
	}
	return nil
}

func (st *Table) check() error {
	if 0 == len(st.Columns) {
		return errors.New(fmt.Sprintf("table[%s] is empty!", st.Name))
	}

	// check the config of
	if err := st.checkConf(); nil != err {
		return err
	}

	// there can be only one auto column, and it must be defined as a columnKey
	autoIncrement := false
	for name, col := range st.Columns {
		if !col.IsAutoIncrement() {
			continue
		}
		if autoIncrement { // only one
			return errors.New(fmt.Sprintf("table[%s] %s can be only one auto column", st.Name, name))
		}

		found := false
		if keys, ok := st.tblKeys[column.PRI]; ok {
			for _, line := range keys.ColumnVec {
				if line == name {
					found = true
					break
				}
			}
		}
		// must be defined as a columnKey
		if !found {
			return errors.New(fmt.Sprintf("table[%s] column %s auto_increment must be index", st.Name, name))
		}

		autoIncrement = true
	}

	// 如果表不存在，不用做后续的兼容判定
	flag, err := st.hasTable()
	if nil != err {
		return err
	}
	if !flag {
		return nil
	}

	if err := st.checkCompatible(); nil != err {
		return err
	}
	return nil
}
