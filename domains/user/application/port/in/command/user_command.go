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
