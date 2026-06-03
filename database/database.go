package database

import (
	"fmt"

	"github.com/itsvagapov/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	dsn := "host=172.22.70.108 user=postgres password=1234 dbname=postgres port=5432 sslmode=disable TimeZone=UTC"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("Ошибка подключения:", err)
		return nil
	}

	if err := db.AutoMigrate(
		&models.Student{},
		&models.Group{},
		&models.Note{},
	); err != nil {

		fmt.Println("Ошибка миграции:", err)
		return nil
	}

	return db
}
