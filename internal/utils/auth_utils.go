package utils

import (
	"time"

	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(user *models.DBUser, config *config.Config) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = user.Login
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(config.SecretKey))
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetUserID(ctx *gin.Context) uuid.UUID {
	if strUserID, ok := ctx.Get("user_id"); ok {
		return strUserID.(uuid.UUID)
	}
	return uuid.Nil
}
