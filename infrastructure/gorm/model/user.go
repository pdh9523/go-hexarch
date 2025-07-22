package model

type User struct {
	BaseModel
	Username string `gorm:"unique;size:20;not null"`
	Nickname string `gorm:"size:20;not null"`
	Role     string `gorm:"size:16;not null"`
}
