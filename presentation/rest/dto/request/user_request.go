package request

type CheckNicknameAvailabilityRequest struct {
	Nickname string `form:"nickname" binding:"required,min=2,max=20"`
}

type CheckUsernameAvailabilityRequest struct {
	Username string `form:"username" binding:"required,min=2,max=16"`
}

type SignUpRequest struct {
	Username string `form:"username" binding:"required,min=2,max=16"`
	Password string `form:"password" binding:"required,min=8,max=32"`
	Nickname string `form:"nickname" binding:"required,min=2,max=20"`
}

type SignInRequest struct {
	Username string `form:"username" binding:"required,min=2,max=16"`
	Password string `form:"password" binding:"required,min=8,max=32"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `form:"current_password" binding:"required,min=8,max=32"`
	NewPassword     string `form:"new_password" binding:"required,min=8,max=32"`
}

type ChangeNicknameRequest struct {
	Nickname string `form:"nickname" binding:"required,min=2,max=20"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
