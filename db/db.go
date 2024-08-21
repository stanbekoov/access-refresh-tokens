package db

import (
	"database/sql"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	Name         string
	Uid          string `gorm:primaryKey`
	RefreshToken string
	Ip           string
	Email        string
}

var Db *gorm.DB

//функция для иницализации БД
func Init() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	name := os.Getenv("DB_NAME")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprint("host=", host, " user=", user, " dbname=", name, " password=", password)

	sqldb, err := sql.Open("pgx", dsn)

	if err != nil {
		panic(err)
	}

	Db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqldb,
	}), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	if !Db.Migrator().HasTable(User{}) {
		Db.Migrator().CreateTable(User{})
	}
}
