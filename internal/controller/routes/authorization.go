package routes

import (
	"fmt"
	"net/http"

	"github.com/Kotletta-TT/bonus-service/config"
	validate "github.com/Kotletta-TT/bonus-service/internal/controller/validators"
	"github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/repository"
	"github.com/Kotletta-TT/bonus-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func RegisterHandler(repo repository.Repository, config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		usr := validate.ValidateUser(ctx)
		if usr == nil {
			return
		}
		usr.HashPassword()
		if err := repo.AddUser(usr); err != nil {
			errors.HandlerErr(err, ctx)
			return
		}
		token, err := utils.GenerateToken(usr.Login, config.SecretKey)
		if err != nil {
			errors.HandlerErr(err, ctx)
			return
		}
		ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token))
		ctx.Status(http.StatusOK)
	}
}

func LoginHandler(repo repository.Repository, config *config.Config) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		usr := validate.ValidateUser(ctx)
		if usr == nil {
			return
		}
		if err := repo.GetUser(usr); err != nil {
			errors.HandlerErr(err, ctx)
			return
		}
		token, err := utils.GenerateToken(usr.Login, config.SecretKey)
		if err != nil {
			errors.HandlerErr(err, ctx)
			return
		}
		ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token))
		ctx.Status(http.StatusOK)
	}
}