package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/vortexfluc/gobank/internal/gobank/storage"
)

type Server struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Run() error {
	e := echo.New()
	e.Use(middleware.Logger())
	e.HideBanner = true
	s.initAccountEndpoints(e)
	s.initLoginEndpoints(e)

	e.Logger.Fatal(e.Start(":8000"))
	return nil
}

func (s *Server) initAccountEndpoints(e *echo.Echo) {
	e.GET("/account", s.handleGetAccount)
	e.GET("/account/:id", s.handleGetAccountById)
	e.POST("/account", s.handleCreateAccount)
	e.POST("/account/:id/transfer", s.handleTransfer)
	e.PUT("/account/:id", s.handleUpdateAccount)
	e.DELETE("/account/:id", s.handleDeleteAccount)
}

func (s *Server) initLoginEndpoints(e *echo.Echo) {
	e.POST("/login", s.handleLogin)
}
