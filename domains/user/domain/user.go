package domain

type User struct {
	Nickname *Nickname
	Username *Username
	Role     Role
}

func NewUser(nickname *Nickname, username *Username) *User {
	return &User{
		Nickname: nickname,
		Username: username,
		Role:     ROLE_USER,
	}
}
