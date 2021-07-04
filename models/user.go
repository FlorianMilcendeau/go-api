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

func (user *User) PrepareGive() {
	user.Password = ""
}

func GetUserById(userId uint) (User, error) {
	var user User

	if err := DB.First(&user, userId).Error; err != nil {
		return user, errors.New("User not found.")
	}

	user.PrepareGive()

	return user, nil
}

func GetUserByName(username string) (User, error) {
	var user User

	err := DB.Model(User{}).Where("username = ?", username).Take(&user).Error

	if err != nil {
		return User{}, err
	}

	user.PrepareGive()

	return user, nil
}

func LoginCheck(username string, password string) (User, string, error) {
	var err error

	user := User{}

	user, err = GetUserByName(username)

	if err != nil {
		return User{}, "", nil
	}

	err = VerifyPassword(password, user.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return User{}, "", nil
	}

	token, err := token.GenerateToken(user.ID)

	if err != nil {
		return User{}, "", err
	}
	return user, token, nil
}
