package demo

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
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

func InitDB(db Database) (*sql.DB, error) {
	path := strings.Join([]string{db.Username, ":", db.Password, "@tcp(", db.Host, ":", db.Port, ")/", db.Dbname, "?charset=utf8"}, "")
	DB, _ := sql.Open("mysql", path)
	DB.SetConnMaxLifetime(100)
	DB.SetMaxIdleConns(10)
	if err := DB.Ping(); err != nil {
		return nil, err
	}
	return DB, nil
}
