package base

import (
	"database/sql"
	"fmt"
	"github.com/beanstalkd/go-beanstalk"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Yaml struct {
	Mysql      Mysql      `yaml:"mysql"`
	BeanStalkd BeanStalkd `yaml:"beanstalkd"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type BeanStalkd struct {
	Host    string `yaml:"host"`
	Port    string `yaml:"port"`
	NetWork string `yaml:"network"`
}

type Config struct {
	Mysql *sql.DB
	Bean  *beanstalk.Conn
}

var Conf Config

func loadYaml(path string) (*Yaml, error) {
	conf := &Yaml{}
	if file, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(file).Decode(conf)
	}
	return conf, nil
}

func initMysql(mysql Mysql) (*sql.DB, error) {
	path := strings.Join([]string{mysql.Username, ":", mysql.Password, "@tcp(", mysql.Host, ":", mysql.Port, ")/", mysql.Dbname, "?charset=utf8"}, "")
	DB, _ := sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		return nil, err
	}
	return DB, nil
}

func initBeanStalkd(bs BeanStalkd) (*beanstalk.Conn, error) {
	c, err := beanstalk.Dial(bs.NetWork, strings.Join([]string{bs.Host, ":", bs.Port}, ""))
	return c, err
}

func Init(yamlPath string) *Config {
	yaml, err := loadYaml(yamlPath)
	if err != nil {
		fmt.Println(err)
	}
	Db, err := initMysql(yaml.Mysql)
	Bs, err := initBeanStalkd(yaml.BeanStalkd)

	conf := &Config{}
	conf.Mysql = Db
	conf.Bean = Bs
	Conf = *conf
	return conf
}
