package handlers

import (
	"net/http"

	"github.com/Kotletta-TT/bonus-service/internal/app/loyalty/entity"
	"github.com/Kotletta-TT/bonus-service/internal/app/loyalty/usecase"
	"github.com/gin-gonic/gin"
)

type UserHandlers struct {
	uc usecase.UserManagmentUseCase
}

func New(uc usecase.UserManagmentUseCase) *UserHandlers {
	return &UserHandlers{uc: uc}
}

func (uh *UserHandlers) CreateUser(c *gin.Context) {
	var user entity.JSONUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := uh.uc.RegisterUser(c.Request.Context(), user.Login, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Authorization", "Bearer "+token)
	c.Status(http.StatusCreated)
}

func (uh *UserHandlers) LoginUser(c *gin.Context) {
	var user entity.JSONUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := uh.uc.LoginUser(c.Request.Context(), user.Login, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Header("Authorization", "Bearer "+token)
	c.Status(http.StatusOK)
}
