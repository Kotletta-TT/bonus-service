package models

import (
	"html"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type DBUser struct {
	ID         uuid.UUID
	Login      string
	Password   string
	Balance    float64
	Withdrawal float64
}

type ViewUser struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ViewUserBalance struct {
	Current   float64 `json:"current" binding:"required"`
	Withdrawn float64 `json:"withdrawn" binding:"required"`
}

func (user *DBUser) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.Login = html.EscapeString(strings.TrimSpace(user.Login))
	return nil
}

func ConvertDBToViewUser(usr *DBUser) *ViewUser {
	return &ViewUser{Login: usr.Login, Password: usr.Password}
}

func ConvertViewToDBUser(usr *ViewUser) *DBUser {
	return &DBUser{Login: usr.Login, Password: usr.Password}
}
