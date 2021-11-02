package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func TokenCheck(c *gin.Context)  {
	fmt.Print("test")
}