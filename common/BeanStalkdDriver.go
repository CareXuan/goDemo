package common

import (
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	"time"
)

func GetBeanConn() *beanstalk.Conn {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		fmt.Print(err)
	}
	return c
}

func Put(c *beanstalk.Conn) {
	id, err := c.Put([]byte("hello"), 1, 10000, 120*time.Second)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(id)
}
