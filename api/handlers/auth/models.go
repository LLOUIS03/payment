package auth

type CreateTokenRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateTokenResponse struct {
	Token string `json:"token"`
}
