package models

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
    host      = "db" // localhost if run locally
    port      = 5432
    user      = "Rus.drbkv"
    password  = "Aa123456"
    dbname    = "postgres"
    secretKey = "your-secret-key"
)

func ConnectDatabase() {
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        host, port, user, password, dbname)
    fmt.Println("DSN: ", dsn)

    var err error
    retryCount := 5 // Количество попыток переподключения
    for i := 0; i < retryCount; i++ {
        DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            break // Если подключение успешно, выходим из цикла
        }
        fmt.Printf("Failed to connect to database (attempt %d): %s\n", i+1, err.Error())
        time.Sleep(time.Second * 5) // Пауза перед следующей попыткой
    }

    if err != nil {
        panic("Failed to connect to database after several attempts")
    }

    DB.AutoMigrate(&Task{}, &User{})
}
