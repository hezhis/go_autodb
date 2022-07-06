package go_autodb

import (
	"fmt"
)

type ProcedureSt struct {
	Name string
	SQL  string
}

func (st *ProcedureSt) Build() error {
	if err := execSQL(st.CreateDropSQL(), false); nil != err {
		return err
	}
	if err := execSQL(st.CreateProcedureSQL(), false); nil != err {
		return err
	}
	return nil
}

func (st *ProcedureSt) CreateDropSQL() string {
	return fmt.Sprintf(ProcedureDropTemplate, st.Name)
}

func (st *ProcedureSt) CreateProcedureSQL() string {
	if len(st.SQL) <= 0 {
		return ""
	}
	return st.SQL
}
