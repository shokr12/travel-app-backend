package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func ConnectToDB() *gorm.DB {
	godotenv.Load()

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_NAME")
	err1 := godotenv.Load(".env")
	if err1 != nil {
		fmt.Println("Error loading .env file")
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, database)
	fmt.Println(dsn)
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return Db
}
