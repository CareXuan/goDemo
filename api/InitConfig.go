package api

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Config struct {
	DB Database `yaml:"database"`
}

type Database struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Dbname   string `yaml:"dbname"`
}

func loadYaml(path string) (*Config, error) {
	conf := &Config{}
	if file, err := os.Open(path); err != nil {
		return nil, err
	} else {
		yaml.NewDecoder(file).Decode(conf)
	}
	return conf, nil
}

func initDB(db Database) (*sql.DB, error) {
	path := strings.Join([]string{db.Username, ":", db.Password, "@tcp(", db.Host, ":", db.Port, ")/", db.Dbname, "?charset=utf8"}, "")
	DB, _ := sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		return nil, err
	}
	return DB, nil
}

func Init(yamlPath string) {
	conf, err := loadYaml(yamlPath)
	if err != nil {
		fmt.Println(err)
	}
	Db, err := initDB(conf.DB)
	fmt.Println(Db)
}
