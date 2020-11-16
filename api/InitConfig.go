package api

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Yaml struct {
	Mysql Mysql `yaml:"mysql"`
}

type Mysql struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

type Config struct {
	Mysql *sql.DB
}

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

func Init(yamlPath string) *Config {
	yaml, err := loadYaml(yamlPath)
	if err != nil {
		fmt.Println(err)
	}
	Db, err := initMysql(yaml.Mysql)

	conf := &Config{}
	conf.Mysql = Db
	return conf
}
