package account

import (
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

type Account struct {
	ID                int
	FirstName         string
	LastName          string
	Number            int64
	EncryptedPassword string
	Balance           int64
	CreatedAt         time.Time
}

func (a *Account) ValidatePassword(pw string) error {
	return bcrypt.CompareHashAndPassword([]byte(a.EncryptedPassword), []byte(pw))
}

func NewAccount(firstName, lastName, password string) (*Account, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &Account{
		FirstName:         firstName,
		LastName:          lastName,
		Number:            int64(rand.Intn(1000000)),
		EncryptedPassword: string(encpw),
		CreatedAt:         time.Now().UTC(),
	}, nil
}
