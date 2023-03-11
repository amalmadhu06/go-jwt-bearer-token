package db

import (
	"fmt"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/config"
	"github.com/amalmadhu06/go-jwt-bearer-token/pkg/domain"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(cfg config.Config) (*gorm.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s user=%s dbname=%s port=%s password=%s", cfg.DBHost, cfg.DBUser, cfg.DBName, cfg.DBPort, cfg.DBPassword)
	db, dbErr := gorm.Open(postgres.Open(psqlInfo), &gorm.Config{
		SkipDefaultTransaction: true,
	})

	err := db.AutoMigrate(

		//add the tables to migrate here
		&domain.Admin{},
		&domain.User{},
		&domain.Product{},
	)
	if err != nil {
		return nil, err
	}

	return db, dbErr
}
