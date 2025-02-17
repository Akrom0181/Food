package models

type UserLoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Phone        string `json:"phone"`
	Id           string `json:"id"`
}

type AuthInfo struct {
	UserID   string `json:"user_id"`
	UserRole string `json:"user_role"`
}

type UserRegisterRequest struct {
	Email string `json:"email"`
}

type UserRegisterConfRequest struct {
	MobilePhone string `json:"email"`
	Otp         string `json:"otp"`
	User        *User  `json:"user"`
}

type UserLoginPhoneConfirmRequest struct {
	Email   string `json:"email"`
	SmsCode string `json:"smscode"`
}
