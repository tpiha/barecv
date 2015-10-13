package main

import (
	"github.com/jinzhu/gorm"
)

// Initializes database
func InitDb() {
	db, _ := gorm.Open("postgres", "user=postgres password=ouy2kcuF  dbname=tinycv sslmode=disable")

	db.DB()

	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)

	db.AutoMigrate(&User{})
}
