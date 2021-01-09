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
	id, err := c.Put([]byte("hello"), 1, 1, 120*time.Second)
	if err != nil {
		fmt.Print(err)
	}
	fmt.Print(id)
}

func Get(c *beanstalk.Conn) {
	c, err := beanstalk.Dial("tcp", "127.0.0.1:11300")
	if err != nil {
		fmt.Print(err)
	}
	id, body, err := c.Reserve(1000 * time.Second)
	fmt.Print(id)
	fmt.Print(string(body))
	//c.Delete(id)
}
