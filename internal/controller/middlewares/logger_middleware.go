package middlewares

import (
	"fmt"
	"time"

	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/gin-gonic/gin"
)

func RequestResponseLogging() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		// for header, value := range ctx.Request.Header {
		// 	fmt.Printf("%s: %s\n", header, value)
		// }
		ctx.Next()
		duration := time.Since(startTime)
		WriteSize := ctx.Writer.Size()
		if WriteSize < 0 {
			WriteSize = 0
		}
		logger.Info(fmt.Sprintf("%s %s %d %s %s %d bytes", ctx.ClientIP(), ctx.Request.Method, ctx.Writer.Status(), ctx.Request.URL.Path, duration, WriteSize))
	}
}
