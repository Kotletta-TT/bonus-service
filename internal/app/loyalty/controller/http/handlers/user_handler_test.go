package handlers

import (
	"testing"

	"github.com/Kotletta-TT/bonus-service/internal/app/loyalty/usecase"
	"github.com/gin-gonic/gin"
)

func TestUserHandlers_LoginUser(t *testing.T) {
	type fields struct {
		uc usecase.UserManagmentUseCase
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uh := &UserHandlers{
				uc: tt.fields.uc,
			}
			uh.LoginUser(tt.args.c)
		})
	}
}
