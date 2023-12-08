package routes

import (
	"net/http"

	"github.com/Kotletta-TT/bonus-service/internal/controller/validators"
	"github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/Kotletta-TT/bonus-service/internal/repository"
	"github.com/Kotletta-TT/bonus-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func OrderSetHandler(repo repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		order := validators.ValidateSetOrder(ctx)
		if order == nil {
			return
		}
		err := repo.SetUserOrder(order)
		if err != nil {
			errors.HandlerErr(err, ctx)
			return
		}
		ctx.Writer.WriteHeader(http.StatusAccepted)
	}
}

func OrdersListHandler(repo repository.Repository) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		usrID := utils.GetUserID(ctx)
		dbOrders, err := repo.GetUserOrders(usrID)
		if err != nil {
			errors.HandlerErr(err, ctx)
			return
		}
		if len(dbOrders) == 0 {
			ctx.Writer.WriteHeader(http.StatusNoContent)
			return
		}
		orders := models.ConvertDBToView(dbOrders)
		ctx.JSONP(http.StatusOK, orders)
	}
}
