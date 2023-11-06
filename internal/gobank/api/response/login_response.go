package response

type LoginResponse struct {
	Number int64  `json:"number"`
	Token  string `json:"jwt_token"`
}
