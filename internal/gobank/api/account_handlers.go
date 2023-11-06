package api

import (
	"github.com/labstack/echo/v4"
	"github.com/vortexfluc/gobank/internal/gobank/account"
	"github.com/vortexfluc/gobank/internal/gobank/api/request"
	"github.com/vortexfluc/gobank/internal/gobank/api/response"
	"net/http"
	"strconv"
)

func (s *Server) handleGetAccount(c echo.Context) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	accountsDto := make([]*response.AccountDto, 0)
	for _, a := range accounts {
		accountsDto = append(accountsDto, response.CreateAccountResponse(a))
	}

	return c.JSON(http.StatusOK, accountsDto)
}

func (s *Server) handleGetAccountById(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return err
	}

	a, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.CreateAccountResponse(a))
}

func (s *Server) handleCreateAccount(c echo.Context) error {
	var req request.CreateAccountRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	a, err := account.NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if _, err := s.store.CreateAccount(a); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, response.CreateAccountResponse(a))
}

func (s *Server) handleDeleteAccount(c echo.Context) error {
	id, err := strconv.Atoi(c.QueryParam("id"))
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]int{"deleted": id})
}

func (s *Server) handleTransfer(c echo.Context) error {
	var transferReq request.TransferRequest
	if err := c.Bind(&transferReq); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, transferReq)
}
