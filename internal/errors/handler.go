package errors

import (
	"errors"
	"net/http"

	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/gin-gonic/gin"
)

func HandlerErr(err error, ctx *gin.Context) {
	logger.Info(err.Error())
	var regErr *UsersError
	if errors.As(err, &regErr) {
		ctx.JSON(regErr.Code, gin.H{"error": regErr.Error()})
		return
	}
	ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
}
