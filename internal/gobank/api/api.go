package api

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/vortexfluc/gobank/internal/gobank/account"
	"github.com/vortexfluc/gobank/internal/gobank/api/request"
	"github.com/vortexfluc/gobank/internal/gobank/api/response"
	jwtutils "github.com/vortexfluc/gobank/internal/gobank/jwt-utils"
	"github.com/vortexfluc/gobank/internal/gobank/storage"
	"log"
	"net/http"
	"strconv"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() error {
	router := mux.NewRouter()

	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleAccountId), s.store))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

	log.Println("JSON API server running on port: ", s.listenAddr)

	err := http.ListenAndServe(s.listenAddr, router)
	if err != nil {
		return err
	}

	return nil
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleAccountId(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetAccountById(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

func (s *APIServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "POST" {
		return WriteJSON(w, http.StatusMethodNotAllowed, APIError{Error: "method not supported"})
	}

	req := new(request.LoginRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	acc, err := s.store.GetAccountByNumber(req.Number)
	if err != nil {
		return err
	}

	err = acc.ValidatePassword(req.Password)
	if err != nil {
		return WriteJSON(w, http.StatusForbidden, APIError{Error: "not authenticated"})
	}

	tokenStr, err := jwtutils.CreateJWT(acc)
	if err != nil {
		return err
	}

	resp := response.LoginResponse{
		Number: acc.Number,
		Token:  tokenStr,
	}

	return WriteJSON(w, http.StatusOK, resp)
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAccounts()
	if err != nil {
		return err
	}

	accountsDto := make([]*response.AccountDto, 0)
	for _, a := range accounts {
		accountsDto = append(accountsDto, response.CreateAccountResponse(a))
	}

	return WriteJSON(w, http.StatusOK, accountsDto)
}

func (s *APIServer) handleGetAccountById(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	a, err := s.store.GetAccountById(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, response.CreateAccountResponse(a))
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	req := new(request.CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return err
	}

	a, err := account.NewAccount(req.FirstName, req.LastName, req.Password)
	if err != nil {
		return err
	}

	if _, err := s.store.CreateAccount(a); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, response.CreateAccountResponse(a))
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := getID(r)
	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	transferReq := new(request.TransferRequest)
	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()

	return WriteJSON(w, http.StatusOK, transferReq)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func permissionDenied(w http.ResponseWriter) {
	WriteJSON(w, http.StatusForbidden, APIError{Error: "permission denied"})
}

func withJWTAuth(handlerFunc http.HandlerFunc, s storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("x-jwt-token")

		token, err := jwtutils.ValidateJWT(tokenStr)
		if err != nil {
			permissionDenied(w)
			log.Println(err)
			return
		}

		if !token.Valid {
			permissionDenied(w)
			return
		}

		userId, err := getID(r)
		if err != nil {
			permissionDenied(w)
			log.Println(err)
			return
		}

		a, err := s.GetAccountById(userId)
		if err != nil {
			permissionDenied(w)
			log.Println(err)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if a.Number != int64(claims["AccountNumber"].(float64)) {
			permissionDenied(w)
			log.Println(a.Number, "!=", claims["AccountNumber"])
			return
		}

		handlerFunc(w, r)
	}
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type APIError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, APIError{Error: err.Error()})
		}
	}
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return id, fmt.Errorf("invalid id given: %s", idStr)
	}

	return id, nil
}
