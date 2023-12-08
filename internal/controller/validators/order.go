package validators

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/Kotletta-TT/bonus-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func ValidateSetOrder(ctx *gin.Context) *models.DBOrders {
	if ctx.GetHeader("Content-Type") != "text/plain" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request type"})
		return nil
	}
	rawBody, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		logger.Info(err.Error())
		return nil
	}
	number := string(rawBody)
	if !utils.LuhnValid(number) {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Invalid order"})
		logger.Debug(fmt.Sprintf("luhn not valid '%s'", number))
		return nil
	}
	usrID := utils.GetUserID(ctx)
	return &models.DBOrders{Number: number, UserID: usrID}
}
