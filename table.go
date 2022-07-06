package go_autodb

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hezhis/go_autodb/column"
	"github.com/hezhis/go_autodb/db"
)

var tables = make(map[string]*Table)

type MysqlTable struct {
	Columns map[string]*column.MysqlColumn
	Keys    map[string]*TableKey
}

type Table struct {
	Name        string
	Comment     string
	Columns     map[string]column.IColumn
	Sequence    []column.IColumn
	Exists      bool
	tblKeys     map[string]*TableKey
	sqlInfo     *MysqlTable
	hasDataFlag int8
}

func NewTable(name, comment string, fieldCount int) *Table {
	table := &Table{
		Name: name, Comment: comment, hasDataFlag: -1,
	}
	table.Sequence = make([]column.IColumn, 0, fieldCount)
	table.Columns = make(map[string]column.IColumn)
	table.tblKeys = make(map[string]*TableKey)
	return table
}

func (st *Table) addColumn(column column.IColumn) error {
	name := column.GetName()
	if _, exists := st.Columns[name]; exists {
		return errors.New(fmt.Sprintf("the column [%s] of [%s] is exists!", name, st.Name))
	}
	st.Columns[name] = column
	st.Sequence = append(st.Sequence, column)
	return nil
}

func (st *Table) createTableSQL() string {
	sqlVec := make([]string, 0, len(st.Sequence))
	for _, column := range st.Sequence {
		sqlVec = append(sqlVec, column.CreateColumnSQL())
	}

	for _, key := range st.tblKeys {
		sql := key.createKeySQL()
		if len(sql) > 0 {
			sqlVec = append(sqlVec, sql)
		}
	}

	head := fmt.Sprintf(CreateSQLHead, st.Name)
	tail := fmt.Sprintf(CreateSQLTail, st.Comment)
	return head + strings.Join(sqlVec, ",\n") + tail
}

func (st *Table) hasTable() (bool, error) {
	if st.Exists {
		return true, nil
	}
	var ret []string
	err := db.OrmEngine.SQL(fmt.Sprintf("show tables like '%s';", st.Name)).Find(&ret)
	if nil != err {
		return false, err
	}

	st.Exists = len(ret) > 0
	return st.Exists, nil
}

func (st *Table) hasData() (bool, error) {
	if st.hasDataFlag != -1 {
		return st.hasDataFlag == 1, nil
	}
	flag, err := st.hasTable()
	if nil != err {
		return false, err
	}
	if !flag {
		st.hasDataFlag = 0
		return false, nil
	}
	rets, err := db.OrmEngine.QueryString(fmt.Sprintf("select * from %s limit 1;", st.Name))
	if nil != err {
		return false, err
	}
	if len(rets) > 0 {
		st.hasDataFlag = 1
	} else {
		st.hasDataFlag = 0
	}
	return st.hasDataFlag == 1, nil
}

func (st *Table) Build() error {
	if err := st.check(); nil != err {
		return err
	}

	flag, err := st.hasTable()
	if nil != err {
		return err
	}
	if !flag {
		st.create()
	} else {
		if err := st.change(); nil != err {
			return err
		}
	}
	return nil
}

func (st *Table) create() {
	execSQL(st.createTableSQL(), true)
}

func (st *Table) change() error {
	bChange := false
	info, err := st.GetTableColumnInfo(false)
	if nil != err {
		return err
	}
	if nil == info {
		return nil
	}
	columns := info.Columns

	for name, column := range st.Columns {
		//该列已经存在，检查是否需要修改
		if mysqlInfo, ok := columns[name]; ok {
			if !column.IsEqual(mysqlInfo) {
				execSQL(column.ChangeColumnSQL(st.Name), true)
			}
		} else { //不存在则添加
			execSQL(column.AddColumnSQL(st.Name), true)
			bChange = true
		}
	}
	//需要删除的column
	for _, mysqlInfo := range columns {
		if _, ok := st.Columns[mysqlInfo.Field]; !ok {
			if err := st.execDropColumnSQL(mysqlInfo.Field); nil != err {
				return err
			}
			bChange = true
		}
	}

	// 表变化了重新查询一次数据库表信息
	if bChange {
		info, err = st.GetTableColumnInfo(true)
		if nil != err {
			return err
		}
		columns = info.Columns
	}
	// key的添加放在删除后，否则可能会冲突
	// 已经检查过key的合法性，所以对不存在的key直接添加，这里不处理 PRIMARY KEY之前只有一个，现在有两个的问题
	for _, key := range st.tblKeys {
		if _, exist := info.Keys[key.Name]; !exist {
			if err := execSQL(key.addKeySQL(st.Name), true); nil != err {
				return err
			}
		}
	}

	//当添加或删除过字段时,认为字段顺序可能不一致
	//比对顺序太麻烦了，直接全部change一遍
	if bChange {
		for i, column := range st.Sequence {
			if i == 0 {
				column.SetFirst()
			} else {
				column.SetPlace(st.Sequence[i-1].GetName())
			}
			if err := execSQL(column.ChangeColumnSQL(st.Name), true); nil != err {
				return err
			}
		}
	}
	return nil
}

func (st *Table) GetTableColumnInfo(reset bool) (*MysqlTable, error) {
	if nil != st.sqlInfo && !reset {
		return st.sqlInfo, nil
	}
	st.sqlInfo = new(MysqlTable)
	st.sqlInfo.Columns = make(map[string]*column.MysqlColumn)

	columns := make([]*column.MysqlColumn, 0)
	if err := db.OrmEngine.SQL(fmt.Sprintf("show columns from %s;", st.Name)).Find(&columns); nil != err {
		return nil, err
	}

	for _, ret := range columns {
		st.sqlInfo.Columns[ret.Field] = ret
	}

	// keys
	vec := make([]*MysqlKeyDesc, 0)
	db.OrmEngine.SQL(fmt.Sprintf("show index from %s;", st.Name)).Find(&vec)
	sort.Slice(vec, func(i, j int) bool {
		return vec[i].SeqInIndex < vec[j].SeqInIndex
	})
	st.sqlInfo.Keys = make(map[string]*TableKey)
	for _, line := range vec {
		if line.KeyName == "PRIMARY" {
			line.KeyName = column.PRI
		}
		cur, exist := st.sqlInfo.Keys[line.KeyName]
		if !exist {
			cur = createTableKeyByMysqlData(line)
			st.sqlInfo.Keys[line.KeyName] = cur
		} else {
			cur.ColumnVec = append(cur.ColumnVec, line.ColumnName)
		}
	}

	return st.sqlInfo, nil
}

func (st *Table) execDropColumnSQL(columnName string) error {
	return execSQL(fmt.Sprintf("alter table %s drop column %s;", st.Name, columnName), true)
}

func (st *Table) DropKeySQL(columnName, kt string) error {
	var sql string
	switch kt {
	case column.MUL:
		sql = fmt.Sprintf("alter table %s drop index %s", st.Name, columnName)
	case column.PRI:
		sql = fmt.Sprintf("alter table %s drop primary columnkey", st.Name)
	case column.UNI:
		sql = fmt.Sprintf("alter table %s drop index %s", st.Name, columnName)
	default:
		return errors.New(fmt.Sprintf("drop key error. table:%s, colName:%s, key type:%s", st.Name, columnName, kt))
	}
	return execSQL(sql, true)
}
