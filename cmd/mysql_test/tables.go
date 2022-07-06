package main

import "github.com/hezhis/go_autodb"

func init() {
	go_autodb.Tables = []*go_autodb.TableConf{
		{
			Name:    "account",
			Comment: "账号表",
			Columns: []*go_autodb.ColumnConf{
				{Name: "user_id", Type: "int", Comment: "用户唯一ID", Unsigned: true},
				{Name: "account_name", Type: "varchar", Comment: "用户帐户的字符串", Size: 64},
				{Name: "passwd", Type: "varchar", Comment: "玩家的密码", Size: 32},
			},
			Keys: []*go_autodb.TableKey{
				{Type: "pri", Name: "user_id"},
				{Type: "uni", Name: "account_name"},
			},
		},
		{
			Name:    "constvariables",
			Comment: "常量表",
			Columns: []*go_autodb.ColumnConf{
				{Name: "actorid_series_bits", Type: "int", Comment: "玩家序列号所占bits，需要同时修改以下掩码"},
				{Name: "actorid_series_mask", Type: "bigint", Comment: "玩家序列号掩码"},
				{Name: "actor_mail_max", Type: "int", Comment: "玩家邮件数量达到最大值后开始替换旧邮件"},
			},
		},
	}
}
