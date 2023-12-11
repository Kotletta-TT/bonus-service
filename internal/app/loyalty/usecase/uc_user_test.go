package usecase_test

import (
	"context"
	"testing"

	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/app/loyalty/entity"
	"github.com/Kotletta-TT/bonus-service/internal/app/loyalty/usecase"
	mock_usecase "github.com/Kotletta-TT/bonus-service/internal/app/loyalty/usecase/mocks"
	mock_utils "github.com/Kotletta-TT/bonus-service/internal/app/loyalty/utils/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestUserUseCase_LoginUser(t *testing.T) {
	type args struct {
		user   *entity.User
		config *config.Config
	}
	tests := []struct {
		name    string
		args    args
		mockErr error
	}{
		{name: "Normal login user, return token",
			args: args{user: &entity.User{
				Login:    "test",
				Password: "test",
			}, config: &config.Config{
				SecretKey: "secret",
			},
			},
			mockErr: nil,
		},
		{name: "User does not exist"},
		{name: "Incorrect verify password"},
		{name: "Invalid genarate token"},
	}
	// DB Error does not exist
	// Hash/Verify password
	// Generate/Validate token
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_usecase.NewMockUserRepo(ctrl)
	mockHasherPasswords := mock_utils.NewMockPassworder(ctrl)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetUserByLogin := mockUserRepo.EXPECT().GetUserByLogin(context.Background(), tt.args.user.Login)
			GetUserByLogin.Return(tt.args.user, nil)
			VerifyPassword := mockHasherPasswords.EXPECT().Verify(tt.args.user.Password, tt.args.user.Password)
			VerifyPassword.Return(nil)
			uc := usecase.NewUserUseCase(mockUserRepo, tt.args.config, mockHasherPasswords)
			token, err := uc.LoginUser(context.Background(), tt.args.user.Login, tt.args.user.Password)
			if tt.mockErr != nil && assert.Error(t, err) {
				assert.Equal(t, tt.mockErr, err)
			}
			// TODO Возможно сюда еще можно добавить проверку токена.
			assert.NotEqual(t, "", token)
		})
	}
}
