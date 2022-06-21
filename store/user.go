package store

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"not null;default:'empty'"`
	Email    string `gorm:"not null;default:'empty';unique"`
	//	EmailVerified bool   `gorm:"not null"`
	Password string `gorm:"not null;default:'empty';"`
	//	Role          string `gorm:"not null;"`
	Picture string
}

func (u *User) Create() (err error) {
	err = db.Save(&u).Error
	return
}

func (u *User) Update(user *User) (err error) {
	err = db.Model(&u).Updates(&user).Error
	return
}

func (u *User) Delete() (err error) {
	err = db.Delete(&u).Error
	return
}

func (u *User) GetUser() (user User, err error) {
	err = db.Where(&u).Find(&user).Error
	return
}

func (u *User) CheckIfExists() (exists bool, err error) {
	exists = false
	err = db.Model(&User{}).Select("count(*) > 0").Where("email = ?", u.Email).Find(&exists).Error
	if err != nil {
		return
	}
	return
}
