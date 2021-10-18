package api

import (
	"github.com/gin-gonic/gin"
)

func Test1(c *gin.Context) {
	GetOk(c, "test ok", []string{})
}
