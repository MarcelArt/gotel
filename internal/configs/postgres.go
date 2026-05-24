package configs

import (
	"fmt"
	"log"
	"strconv"

	"github.com/MarcelArt/gotel/internal/entities"
	"github.com/MarcelArt/gotel/internal/enums"
	"github.com/MarcelArt/gotel/pkg/jsonb"
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
		entities.Role{},
		entities.UserRole{},
		entities.Category{},
		entities.Item{},
		entities.Location{},
		entities.InventoryTransaction{},
	)
	fmt.Println("Database Migrated")

	seedDefaultUser()
	fmt.Println("Default User Seeded")
	seedLocations()
	fmt.Println("Default Locations Seeded")
	seedCategories()
	fmt.Println("Default Categories Seeded")

	return err
}

func DropDB() error {
	db := DB
	err := db.Migrator().DropTable(
		entities.User{},
		entities.Role{},
		entities.UserRole{},
		entities.Category{},
		entities.Item{},
		entities.Location{},
		entities.InventoryTransaction{},
	)
	fmt.Println("Database Dropped")

	return err
}

func seedDefaultUser() {
	user := entities.User{
		Username: Env.DefaultUser,
		Email:    Env.DefaultEmail,
		Password: Env.DefaultPassword,
	}
	DB.Where("username = ?", user.Username).FirstOrCreate(&user)

	permissions, err := jsonb.New([]string{enums.PermFullAccess})
	if err != nil {
		log.Fatalf("failed seeding role: %s", err.Error())
	}

	role := entities.Role{
		Name:        "Admin",
		Permissions: permissions,
	}
	DB.Where("name = ?", role.Name).FirstOrCreate(&role)

	userRole := entities.UserRole{
		UserID: user.ID,
		RoleID: role.ID,
	}
	DB.Where("role_id = ? and user_id = ?", role.ID, user.ID).FirstOrCreate(&userRole)
}

func seedLocations() {
	locations := []entities.Location{
		{
			Value:       "Main Warehouse",
			IsVirtual:   false,
			Description: "Primary inventory storage",
		},
		{
			Value:       "Laundry",
			IsVirtual:   false,
			Description: "Laundry processing area",
		},
		{
			Value:       "Disposed",
			IsVirtual:   true,
			Description: "For disposed or damaged items",
		},
		{
			Value:       "Missing",
			IsVirtual:   true,
			Description: "For lost or unaccounted items",
		},
		{
			Value:       "Consumed",
			IsVirtual:   true,
			Description: "For consumable items already used",
		},
	}

	for _, location := range locations {
		DB.Where("value = ?", location.Value).FirstOrCreate(&location)
	}
}

func seedCategories() {
	categories := []entities.Category{
		{Value: "Linen"},
		{Value: "Amenities"},
		{Value: "Cleaning Supplies"},
		{Value: "Equipment"},
	}

	for _, category := range categories {
		DB.Where("value = ?", category.Value).FirstOrCreate(&category)
	}
}
