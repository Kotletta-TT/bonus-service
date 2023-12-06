package validators

import (
	"net/http"

	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/gin-gonic/gin"
)

func ValidateUser(ctx *gin.Context) *models.DBUser {
	var requestUser models.ViewUser
	if err := ctx.ShouldBindJSON(&requestUser); err != nil {
		logger.Info(err.Error())
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil
	}
	return models.ConvertViewToDBUser(&requestUser)
}
