package utils

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func GenerateToken(login string, secretKey string) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["id"] = login
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(1)).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secretKey))
}

type Passworder interface {
	Hash(password string) (string, error)
	Verify(password, hashedPassword string) error
}

type HasherPassworder struct {
}

func NewHashPassworder() *HasherPassworder {
	return &HasherPassworder{}
}

func (hp *HasherPassworder) Hash(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}
	password = string(hashedPassword)
	// user.Login = html.EscapeString(strings.TrimSpace(user.Login))
	return password, nil
}

func (hp *HasherPassworder) Verify(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func GetUserID(ctx *gin.Context) uuid.UUID {
	if strUserID, ok := ctx.Get("user_id"); ok {
		return strUserID.(uuid.UUID)
	}
	return uuid.Nil
}
