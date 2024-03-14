package user

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type UpdateRequest struct {
	Nickname *string `json:"nickname"`
	Email    *string `json:"email"`
}
