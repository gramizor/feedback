package main

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"rest-apishka/internal/dsn"
	"rest-apishka/internal/model"
)

func main() {
	_ = godotenv.Load()
	db, err := gorm.Open(postgres.Open(dsn.FromEnv()), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Явно мигрировать только нужные таблицы
	err = db.AutoMigrate(&model.Group{}, &model.Feedback{}, &model.User{}, &model.FeedbackGroup{})
	if err != nil {
		panic("cant migrate db")
	}
}
