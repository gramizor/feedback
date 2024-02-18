package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"rest-apishka/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const numRecords = 500000

func main() {
	dsn := "host=127.0.0.1 port=5434 user=postgres password=123 dbname=feedback sslmode=disable"
	// Инициализация подключения к базе данных
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	// Автомиграция для создания таблицы groups, если ее нет
	db.AutoMigrate(&model.Group{})

	// Загрузка 10 тысяч записей
	err = insertRandomGroups(db, numRecords)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%d записей успешно загружены в таблицу groups\n", numRecords)
}

// Вставка случайных записей в таблицу groups
func insertRandomGroups(db *gorm.DB, numRecords int) error {
	for i := 1; i <= numRecords; i++ {

		group := model.Group{
			GroupCode:   generateGroupCode(i),
			Contacts:    fmt.Sprintf("+7(999) 999-99-9%d", i%10),
			Course:      generateCourse(i),
			Students:    rand.Intn(29) + 2,
			GroupStatus: model.GROUP_STATUS_ACTIVE,
			Photo:       "",
		}

		// Начало транзакции
		tx := db.Begin()

		if err := tx.Create(&group).Error; err != nil {
			// Откат транзакции при ошибке
			tx.Rollback()
			return err
		}

		// Фиксация транзакции
		tx.Commit()

		// Небольшая задержка для имитации реального использования
		time.Sleep(time.Millisecond)
	}

	return nil
}

func generateGroupCode(i int) string {
	rand.Seed(time.Now().UnixNano())
	// Выбор рандомно между ИУ и РТ
	code := "ИУ"
	if rand.Intn(2) == 0 {
		code = "РТ"
	}
	// Цифра рандомная от 1 до 5
	digit1 := rand.Intn(5) + 1
	// Цифра, отображающая номер семестра (от 1 до 8)
	digit2 := i%8 + 1
	// Цифра, отображающая номер группы в потоке (от 1 до 5)
	digit3 := rand.Intn(5) + 1

	return fmt.Sprintf("%s%d-%d%dБ", code, digit1, digit2, digit3)
}

func generateCourse(i int) int {
	digit2 := i%8 + 1
	switch digit2 {
	case 1, 2:
		return 1
	case 3, 4:
		return 2
	case 5, 6:
		return 3
	case 7, 8:
		return 4
	default:
		return 1
	}
}
