package request

type LoginRequest struct {
	Number   int64  `json:"number"`
	Password string `json:"password"`
}
