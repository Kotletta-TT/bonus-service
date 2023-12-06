package middlewares

import (
	"fmt"
	"strings"

	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/Kotletta-TT/bonus-service/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func getRequestToken(ctx *gin.Context) (string, error) {
	authHeader := ctx.GetHeader("Authorization")
	logger.Debug(fmt.Sprintf("get auth header: %s", authHeader))
	splitHeader := strings.Split(authHeader, " ")
	if len(splitHeader) == 2 {
		logger.Debug(fmt.Sprintf("split auth header: %s", splitHeader[1]))
		return splitHeader[1], nil
	}
	return "", errors.InvalidTokenErr()
}

func GetToken(ctx *gin.Context, config *config.Config) (*jwt.Token, error) {
	stringToken, err := getRequestToken(ctx)
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse(stringToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logger.Info(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
			return nil, errors.InvalidTokenErr()
		}
		return []byte(config.SecretKey), nil
	})
	if err != nil {
		logger.Debug(err.Error())
		err = errors.ExpiredTokenErr()
	}
	return token, err
}

func ValidateToken(ctx *gin.Context, config *config.Config) (string, error) {
	token, err := GetToken(ctx, config)

	if err != nil {
		return "", err
	}

	validTokenClaims, ok := token.Claims.(jwt.MapClaims)
	logger.Debug(fmt.Sprintf("parse token id:'%s' auth:'%v' exp:'%.f'", validTokenClaims["id"], validTokenClaims["authorized"], validTokenClaims["exp"]))
	if ok && token.Valid {
		if login, ok := validTokenClaims["id"]; ok {
			return login.(string), nil
		}
	}

	return "", errors.InvalidTokenErr()
}

func Auth(config *config.Config, repo repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		login, err := ValidateToken(ctx, config)
		if err != nil {
			errors.HandlerErr(err, ctx)
			ctx.Abort()
			return
		}
		usr := &models.DBUser{Login: login}
		err = repo.GetTokenUser(usr)
		if err != nil {
			errors.HandlerErr(err, ctx)
			ctx.Abort()
			return
		}
		ctx.Set("user_id", usr.ID)
		ctx.Next()
	}
}
