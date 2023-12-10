package entity

import "github.com/google/uuid"

type User struct {
	ID         uuid.UUID
	Login      string
	Password   string
	Balance    float64
	Withdrawal float64
}

type JSONUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
