package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mouse/base"
)

func Test1(c *gin.Context) {
	bs := base.Conf.Bean
	id, err := base.Put(bs, "carexuan_test", "{\"key\":\"你好呀\"}", 2)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(id)
	//_, body, err := common.Get(bs, "carexuan_test")
	//if err != nil {
	//	fmt.Print(err)
	//}
	//m := mmm{}
	//err = json.Unmarshal(body, &m)
	//if err != nil {
	//	fmt.Print(err)
	//}
	//fmt.Print(m.Key)
	base.GetOk(c, "test ok", []string{})
}
