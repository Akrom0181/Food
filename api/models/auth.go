package models

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AuthInfo struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

type UserRegisterRequest struct {
	Mail string `json:"mail"`
}

type UserRegisterConfRequest struct {
	Mail string `json:"mail"`
	Otp  string `json:"otp"`
	User *User  `json:"user"`
}
