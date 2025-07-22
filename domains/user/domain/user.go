package domain

type User struct {
	ID       string
	Nickname *Nickname
	Username *Username
	Role     Role
}

func NewUser(nickname *Nickname, username *Username) *User {
	return &User{
		ID:       "",
		Nickname: nickname,
		Username: username,
		Role:     ROLE_USER,
	}
}
