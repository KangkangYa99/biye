package pgsql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connectStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	DB, err = sql.Open("postgres", connectStr)
	if err != nil {
		return err
	}
	//限制同时打开的连接数量，防止数据库过载
	DB.SetMaxOpenConns(20)
	//设置最大空闲连接数为 25
	DB.SetMaxIdleConns(10)
	if err = DB.Ping(); err != nil {
		return err
	}
	return nil
}
func CloseDB() {
	if DB != nil {
		if err := DB.Close(); err != nil {

			return
		}
	}
}
