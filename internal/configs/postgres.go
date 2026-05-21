package configs

import (
	"fmt"
	"strconv"

	"github.com/MarcelArt/gotel/internal/entities"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var dsn string

func ConnectDB() {
	p := Env.DBPort
	port, err := strconv.ParseUint(p, 10, 32)
	dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable", Env.DBHost, port, Env.DBUser, Env.DBPassword, Env.DBName, Env.DBSchema)

	if err != nil {
		panic("failed to parse database port")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		panic("failed to connect database")
	}

	DB = db

	fmt.Println("Connection Opened to Database")
}

func MigrateDB() error {
	db := DB
	err := db.AutoMigrate(
		entities.User{},
	)
	fmt.Println("Database Migrated")

	return err
}

func DropDB() error {
	db := DB
	err := db.Migrator().DropTable(
		entities.User{},
	)
	fmt.Println("Database Dropped")

	return err
}
