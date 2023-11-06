package storage

import (
	_ "github.com/lib/pq"
	"github.com/vortexfluc/gobank/internal/gobank/account"
)

type Storage interface {
	CreateAccount(*account.Account) (int, error)
	DeleteAccount(int) error
	UpdateAccount(*account.Account) error
	GetAccounts() ([]*account.Account, error)
	GetAccountById(int) (*account.Account, error)
	GetAccountByNumber(int64) (*account.Account, error)
}
