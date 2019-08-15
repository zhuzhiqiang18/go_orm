package main

import (
	"fmt"
	"github.com/zhuzhiqiang18/go_orm"
)

func main() {
	/*for i:=0;i<10000;i++ {
		tests.TestSave()
	}*/
	//tests.TestSave()
	//tests.TestDelete()
	//tests.TestUpdate()
	//tests.TestFindQuery()
	//tests.TestFindQueryField()
	//tests.TestFindQueryWhere()
	//tests.TestNativeSql()
	//tests.TestAutoInsertId()
	//tests.TestTx()
	//tests.TestTx1()
	var gql go_orm.Gql
	sql:=gql.Where("name = ?").Where("age = ?").Or("is_reading = ?").Order("id desc").Count().GetGql()
	fmt.Println(sql)

	fmt.Println(new(go_orm.Gql).Where("name = ?").Where("age = ?").Or("is_reading = ?").Order("id desc").New().GetGql())

}
