package main

import "github.com/hezhis/go_autodb"

func init() {
	go_autodb.Procedures = []*go_autodb.ProcedureSt{
		{
			Name: "selectchargeorder",
			SQL: `
			create procedure selectchargeorder (in npf_id integer, in nsrv_id integer)
			BEGIN
				select id, actor_id, cash_num, charge_id, pay_no, cp_no from charge where srv_id=nsrv_id and check_time=0;
			END
		`,
		},
		{
			Name: "updateorderchecktime",
			SQL: `
			create procedure updateorderchecktime (in nid bigint)
			BEGIN
				update charge set check_time = UNIX_TIMESTAMP() where id=nid and check_time=0;
			END
		`,
		},
		{
			Name: "dailyclearorder",
			SQL: `
			create procedure dailyclearorder ()
			BEGIN
				delete from charge where check_time != 0 and UNIX_TIMESTAMP()-check_time>=86400;
			END
		`,
		},
	}
}
