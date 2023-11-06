package entity

import "time"

type Account struct {
	ID                int
	FirstName         string
	LastName          string
	Number            int64
	EncryptedPassword string
	Balance           int64
	CreatedAt         time.Time
}
