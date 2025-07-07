package model

import "time"

type User struct {
	Userid   int       `gorm:"primaryKey;not null"`
	Username string    `gorm:"unique;not null" json:"username" form:"username"`
	Password string    `gorm:"not null" json:"password" form:"password"`
	BirthDay time.Time `gorm:"column:birthday;type:date;not null" json:"birthday" form:"birthday"`
}
