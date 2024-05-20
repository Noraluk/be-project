package database

import (
	"be-project/api/entities"
	"be-project/pkg/config"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() error {
	config := config.GetConfig()
	psqlConn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Bangkok", config.Database.Host, config.Database.Username, config.Database.Password, config.Database.DatabaseName, config.Database.Port)

	var err error
	db, err = gorm.Open(postgres.Open(psqlConn), &gorm.Config{})
	if err != nil {
		return err
	}

	db.AutoMigrate(&entities.Auth{}, &entities.Chat{}, &entities.ChatUser{})

	return nil
}

func GetDatabase() *gorm.DB {
	return db
}
