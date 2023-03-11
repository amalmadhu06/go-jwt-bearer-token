package domain

import (
	_ "github.com/go-playground/validator/v10"
	"time"
)

//todo : either use json:exclude or json:"-"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Email     string `gorm:"unique;not null" validate:"email"`
	Password  string `gorm:"not null" validate:"min=8"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
