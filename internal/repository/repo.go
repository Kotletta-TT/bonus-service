package repository

import (
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	GetUser(user *models.DBUser) error
	GetTokenUser(user *models.DBUser) error
	AddUser(user *models.DBUser) error
	GetUserOrders(userID uuid.UUID) ([]*models.DBOrders, error)
	SetUserOrder(order *models.DBOrders) error
	GetUserBalance(userID uuid.UUID) (*models.ViewUserBalance, error)
	GetUserWithdraws(userID uuid.UUID) ([]*models.ViewWithdraw, error)
	RequestUserWithdraw(*models.DBWithdraw) error
	GetUnprocessedOrders() ([]*models.DBOrders, error)
	UpdateOrderAndBalance(*models.DBOrders) error
	Close() error
}
