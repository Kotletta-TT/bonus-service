package routes

import (
	"net/http"

	"github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/repository"
	"github.com/Kotletta-TT/bonus-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func GetBalanceHandler(repo repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := utils.GetUserID(ctx)
		dbBalance, err := repo.GetUserBalance(userID)
		if err != nil {
			errors.HandlerErr(err, ctx)
			return
		}
		ctx.JSON(http.StatusOK, dbBalance)
	}
}
