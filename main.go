package main

import (
	"carexuan/api"
	"carexuan/common"
	"fmt"
)

type T1 struct {
	Id  int64
	Val int64
}

func main() {
	conf := api.Init("./api/config.yaml")
	bs := conf.Bean
	common.Put(bs)
}

func mysqlTest() {
	conf := api.Init("./api/config.yaml")
	mysqlConn := conf.Mysql
	defer mysqlConn.Close()
	sql := "select * from t1"
	rows, err := mysqlConn.Query(sql)
	if err != nil {
		fmt.Print(err)
	}
	var T1s []T1
	for rows.Next() {
		var id int64
		var val int64
		var t1 T1
		rows.Scan(&id, &val)
		t1.Id = id
		t1.Val = val
		T1s = append(T1s, t1)
	}
	fmt.Print(T1s)
}
