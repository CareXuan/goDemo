package main

import (
	api "HelloWorld/api"
	"fmt"
)

type T1 struct {
	Id  int64
	Val string
}

func main() {
	conf := api.Init("./api/config.yaml")
	stmt, err := conf.Mysql.Prepare("select * from t1")
	if err != nil {
		fmt.Println(err)
	}
	rows, err := stmt.Query()

	var T1s []T1
	for rows.Next() {
		var id int64
		var val string
		var t1 T1
		err := rows.Scan(&id, &val)
		if err != nil {
			fmt.Println(err)
		}
		t1.Id = id
		t1.Val = val
		T1s = append(T1s, t1)
	}
	fmt.Println(T1s)
}
