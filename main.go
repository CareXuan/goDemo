package main

import (
	"mouse/common"
)

func main() {
	common.Init("./conf/local.yaml")
	common.InitGin()
}
