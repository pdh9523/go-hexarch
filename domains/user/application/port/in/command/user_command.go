package command

type SignUpCommand struct {
	Username string
	Password string
	Nickname string
}

type SignInCommand struct {
	Username string
	Password string
}

type ChangePasswordCommand struct {
	CurrentPassword string
	NewPassword     string
}

type ChangeNicknameCommand struct {
	Nickname string
}

type RefreshTokenCommand struct {
	RefreshToken string
}
