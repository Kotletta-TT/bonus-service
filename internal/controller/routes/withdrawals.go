package routes

import (
	"net/http"

	"github.com/Kotletta-TT/bonus-service/internal/controller/validators"
	"github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/repository"
	"github.com/Kotletta-TT/bonus-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func RequestWithDrawHandler(repo repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		withdrawDB := validators.ValidateRequestWithdraw(ctx)
		if withdrawDB == nil {
			return
		}
		err := repo.RequestUserWithdraw(withdrawDB)
		if err != nil {
			errors.HandlerErr(err, ctx)
			return
		}
		ctx.Writer.WriteHeader(http.StatusOK)
	}
}

func GetWithDrawalsHandler(repo repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := utils.GetUserID(ctx)
		withdrawsDB, err := repo.GetUserWithdraws(userID)
		if err != nil {
			logger.Error(err.Error())
			errors.HandlerErr(err, ctx)
			return
		}
		// viewViwdrawals := models.ConvertWithdrawDBView(withdrawsDB)
		ctx.JSON(http.StatusOK, withdrawsDB)
	}
}
