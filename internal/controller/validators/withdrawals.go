package validators

import (
	"net/http"

	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/Kotletta-TT/bonus-service/internal/utils"
	"github.com/gin-gonic/gin"
)

func ValidateRequestWithdraw(ctx *gin.Context) *models.DBWithdraw {
	usrID := utils.GetUserID(ctx)
	withdrawRequest := models.WithdrawRequest{}
	if err := ctx.ShouldBindJSON(&withdrawRequest); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		logger.Info(err.Error())
		return nil
	}
	return &models.DBWithdraw{UserID: usrID, OrderID: withdrawRequest.OrderID, Sum: withdrawRequest.Sum}
}
