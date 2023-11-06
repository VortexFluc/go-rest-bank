package storage

import (
	_ "github.com/lib/pq"
	"github.com/vortexfluc/gobank/internal/gobank/types"
)

type Storage interface {
	CreateAccount(*types.Account) error
	DeleteAccount(int) error
	UpdateAccount(*types.Account) error
	GetAccounts() ([]*types.Account, error)
	GetAccountById(int) (*types.Account, error)
	GetAccountByNumber(int64) (*types.Account, error)
}
