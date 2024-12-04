package store

import (
	"os"

	"github.com/kemboy-svg/investment/helpers"
	"github.com/kemboy-svg/investment/helpers/logger"
	"github.com/kemboy-svg/investment/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB

func DbInit() error {
	// envErr := godotenv.Load()
	// if envErr != nil {
	// 	fmt.Printf("Error loading credentials: %v", envErr)
	// }
	dsn := helpers.EliasLocal
	if os.Getenv("IS_PRODUCTION") == "TRUE" {
		dsn = helpers.EliasLocal
	}
	var err error
	Db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		logger.ErrLog(err, "init db error")
		return err
	}

	return nil
}

type Store struct {
	db *gorm.DB
}

func (s Store) MigrateAllModels() {
	m := Db.Migrator()

	m.AutoMigrate(models.Deposit{})

}

func NewStore(db *gorm.DB) *Store {
	return &Store{db: db}
}

