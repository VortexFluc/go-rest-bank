package entity

import (
	"gorm.io/gorm"
)

type Account struct {
	gorm.Model
	ID                int
	FirstName         string
	LastName          string
	Number            int64
	EncryptedPassword string
	Balance           int64
}
