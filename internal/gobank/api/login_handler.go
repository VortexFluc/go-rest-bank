package api

import (
	"github.com/labstack/echo/v4"
	"github.com/vortexfluc/gobank/internal/gobank/api/request"
	"github.com/vortexfluc/gobank/internal/gobank/api/response"
	jwtutils "github.com/vortexfluc/gobank/internal/gobank/jwt-utils"
	"net/http"
)

func (s *Server) handleLogin(c echo.Context) error {
	var req request.LoginRequest

	err := c.Bind(&req)
	if err != nil {
		return c.String(http.StatusBadRequest, "bad request")
	}

	acc, err := s.store.GetAccountByNumber(req.Number)
	if err != nil {
		return c.String(http.StatusForbidden, "bad credentials")
	}

	err = acc.ValidatePassword(req.Password)
	if err != nil {
		return c.String(http.StatusForbidden, "bad credentials")
	}

	tokenStr, err := jwtutils.CreateJWT(acc)
	if err != nil {
		return c.String(http.StatusForbidden, "bad credentials")
	}

	resp := response.LoginResponse{
		Number: acc.Number,
		Token:  tokenStr,
	}

	return c.JSON(http.StatusOK, resp)
}
