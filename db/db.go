package db

import (
	"fmt"
	"log"
	"os"

	"github.com/Anwarjondev/task-management-api/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)



var DB *gorm.DB
func Connect() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error with .env file or not such directory")
	}
	host := os.Getenv("DB_HOST")
	passwrod := os.Getenv("DB_PASSWORD")
	db_user := os.Getenv("DB_USER")
	port := os.Getenv("DB_PORT")
	db_name := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s user=%s password=%s db_name=%s port=%s sslmode=disable TimeZone=UTC",
	host, db_user, passwrod, db_name, port)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		panic("Failed to the connect databse" + err.Error())
	}
	log.Println("Successfully connect to the connect database")
}

func AutoMigrate() {
	err := DB.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Subtask{},
		&models.Task{},
	)
	if err != nil {
		panic("Failed to migrate database")
	}
	log.Println("Database migrated Successfully")
}