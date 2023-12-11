package utils

import (
	"fmt"
	"testing"

	"github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

// func TestVerifyPassword(t *testing.T) {
// 	type args struct {
// 		password       string
// 		hashedPassword string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "Valid password",
// 			args: args{
// 				password:       "password",
// 				hashedPassword: "$2a$10$e28v1ZD5rYz2P4B5guN8U.Cye8pdLT44.Kdq3l1kfNxqCvo2umAdy",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Invalid password",
// 			args: args{
// 				password:       "weak_password",
// 				hashedPassword: "$2a$10$e28v1ZD5rYz2P4B5guN8U.Cye8pdLT44.Kdq3l1kfNxqCvo2umAdy",
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := VerifyPassword(tt.args.password, tt.args.hashedPassword)
// 			if err != nil && tt.wantErr && assert.Error(t, err) {
// 				assert.EqualError(t, err, "crypto/bcrypt: hashedPassword is not the hash of the given password")
// 			}
// 		})
// 	}
// }

func TestGenerateToken(t *testing.T) {
	type args struct {
		login     string
		secretKey string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid login & secret key",
			args: args{
				login:     "user1",
				secretKey: "secret1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := GenerateToken(tt.args.login, tt.args.secretKey)
			if tt.wantErr {
				assert.Error(t, err)
			}
			token, err := jwt.Parse(gotToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					logger.Info(fmt.Sprintf("unexpected signing method: %v", token.Header["alg"]))
					return nil, errors.InvalidTokenErr()
				}
				return []byte(tt.args.secretKey), nil
			})
			if tt.wantErr {
				assert.Error(t, err)
			}
			validTokenClaims, ok := token.Claims.(jwt.MapClaims)
			if ok && token.Valid {
				if login, ok := validTokenClaims["id"]; ok {
					assert.Equal(t, tt.args.login, login)
				} else {
					t.Errorf("no login in token claims")
				}
			} else {
				t.Errorf("token not valid")
			}
		})
	}
}
