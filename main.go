package main

import (
	"carexuan/api"
	"carexuan/common"
	"encoding/json"
	"fmt"
)

type T1 struct {
	Id  int64
	Val int64
}
type mmm struct {
	Key string `json:"key"`
}

func main() {
	conf := api.Init("./conf/local.yaml")
	bs := conf.Bean
	//id, err := common.Put(bs, "carexuan_test", "{\"key\":\"你好呀\"}", 2)
	//if err != nil {
	//	fmt.Print(err)
	//}
	//fmt.Print(id)
	_, body, err := common.Get(bs, "carexuan_test")
	if err != nil {
		fmt.Print(err)
	}
	m := mmm{}
	err = json.Unmarshal(body, &m)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(m.Key)
}

func mysqlTest() {
	conf := api.Init("./conf/local.yaml")
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
