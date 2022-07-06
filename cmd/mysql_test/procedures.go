package main

import "github.com/hezhis/go_autodb"

func init() {
	go_autodb.Procedures = []*go_autodb.ProcedureSt{
		{
			Name: "initdb",
			SQL: `
			create procedure initdb ()
			begin
				delete from constvariables;
				insert into constvariables(actorid_series_bits, actorid_series_mask, actor_mail_max) values(32, 0xFFFFFFFF, 200);
			end
		`,
		},
	}
}
