package main

import (
	"mouse/base"
	"mouse/http"
)

func main() {
	base.Init("./conf/local.yaml")
	http.InitGin()
}
