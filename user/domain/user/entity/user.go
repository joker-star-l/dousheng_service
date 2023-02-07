package entity

type User struct {
	Id       int64  `json:"id" gorm:"column:id"`
	Name     string `json:"name" gorm:"column:name"`
	Password string `json:"password" gorm:"column:password"`
}

func (u *User) TableName() string {
	return "user"
}
