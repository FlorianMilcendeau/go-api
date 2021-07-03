package models

import (
	"api/utils/token"
	"errors"
	"html"
	"strings"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `gorm:"size:255;not null;unique" json:"username"`
	Password string `gorm:"size:255;not null;" json:"password"`
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(username string, password string) (string, error) {
	var err error

	u := User{}

	err = DB.Model(User{}).Where("username = ?", username).Take(&u).Error

	if err != nil {
		return "", nil
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", nil
	}

	token, err := token.GenerateToken(u.ID)

	if err != nil {
		return "", err
	}
	return token, nil
}

func (u *User) SaveUser() (*User, error) {

	var err error

	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}

func (u *User) PrepareGive() {
	u.Password = ""
}

func GetUserById(userId uint) (User, error) {
	var user User

	if err := DB.First(&user, userId).Error; err != nil {
		return user, errors.New("User not found.")
	}

	user.PrepareGive()

	return user, nil
}
