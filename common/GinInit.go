package common

import "github.com/gin-gonic/gin"

func InitGin() {
	r := gin.Default()
	Route(r)
	r.Run()
}
