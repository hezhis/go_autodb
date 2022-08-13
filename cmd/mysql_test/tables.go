package main

import "github.com/hezhis/go_autodb"

func init() {
	go_autodb.Tables = []*go_autodb.TableConf{
		{
			Name:    "charge",
			Comment: "充值订单表",
			Columns: []*go_autodb.ColumnConf{
				{Name: "id", Type: "bigint", Comment: "自增id", Unsigned: true, AutoIncrement: true},
				{Name: "pf_id", Type: "int", Comment: "平台id", Unsigned: true},
				{Name: "srv_id", Type: "int", Comment: "服务器id", Unsigned: true},
				{Name: "account_id", Type: "bigint", Comment: "账号id", Unsigned: true},
				{Name: "actor_id", Type: "bigint", Comment: "角色id", Unsigned: true},
				{Name: "charge_id", Type: "int", Comment: "充值表id", Unsigned: true},
				{Name: "cash_num", Type: "int", Comment: "钱币数量", Unsigned: true},
				{Name: "count", Type: "int", Comment: "购买数量", Unsigned: true},
				{Name: "check_time", Type: "int", Comment: "服务器id", Unsigned: true},
				{Name: "insert_time", Type: "int", Comment: "服务器id", Unsigned: true},
				{Name: "pay_no", Type: "varchar", Size: 64, Comment: "渠道订单号"},
				{Name: "cp_no", Type: "int", Size: 64, Comment: "我们订单号"},
			},
			Keys: []*go_autodb.TableKey{
				{Type: "pri", Name: "pri", ColumnVec: []string{"id"}},
				{Type: "mul", Name: "pf_id", ColumnVec: []string{"pf_id"}},
				{Type: "mul", Name: "srv_id", ColumnVec: []string{"srv_id"}},
				{Type: "mul", Name: "pf_srv_id", ColumnVec: []string{"pf_id", "srv_id"}},
			},
		},
	}
}
