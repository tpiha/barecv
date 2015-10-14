package main

import (
	"github.com/jinzhu/gorm"
	"log"

	_ "github.com/lib/pq"
)

// InitDB initializes database
func InitDb() *gorm.DB {
	log.Printf("[InitDb] pg: %s", config.Postgres)
	db, err := gorm.Open("postgres", config.Postgres)

	if err != nil {
		log.Printf("[InitDb] error: %s", err)
	}

	db.DB()
	db.DB().Ping()
	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.LogMode(config.LogSQL)

	db.AutoMigrate(&User{})

	return &db
}
