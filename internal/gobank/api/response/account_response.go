package response

import (
	"github.com/vortexfluc/gobank/internal/gobank/account"
)

type AccountDto struct {
	ID        int    `json:"id"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Number    int64  `json:"number"`
	Balance   int64  `json:"balance"`
}

func CreateAccountResponse(a *account.Account) *AccountDto {
	return &AccountDto{
		ID:        a.ID,
		FirstName: a.FirstName,
		LastName:  a.LastName,
		Number:    a.Number,
		Balance:   a.Balance,
	}
}
