package postgres

import (
	"errors"
	"fmt"

	customErr "github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/Kotletta-TT/bonus-service/internal/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const insertUser = `
	INSERT INTO users (login, password)
	VALUES ($1, $2)
	RETURNING id
`

const getUser = `
	SELECT id, password, current, withdrawn
	FROM users
	WHERE login = $1 LIMIT 1`

const getUserBalance = `
	SELECT current, withdrawn
	FROM users
	WHERE id = $1 LIMIT 1`

const updateUserBalance = `
	UPDATE users
	SET current = $1, withdrawn = $2
	WHERE id = $3`

func (p *PostgreRepo) GetUser(user *models.DBUser) error {
	return p.WrapConnectionPgError(func() error {
		var dbPassword string
		query := p.pgPool.QueryRow(p.ctx, getUser, user.Login)
		if err := query.Scan(&user.ID, &dbPassword, &user.Balance, &user.Withdrawal); err != nil {
			logger.Debug(err.Error())
			var pgError *pgconn.PgError
			if errors.As(err, &pgError) {
				return err
			}
			return &customErr.UsersError{Code: 401, Err: "unknown login/password"}
		}
		logger.Debug(fmt.Sprintf("get user on db id:'%s' login:'%s'", user.ID.String(), user.Login))
		if err := utils.VerifyPassword(user.Password, dbPassword); err != nil {
			logger.Debug(fmt.Sprintf("user_id: %s verify password err", user.ID))
			return &customErr.UsersError{Code: 401, Err: "unknown login/password"}
		}
		return nil
	})
}

func (p *PostgreRepo) GetTokenUser(user *models.DBUser) error {
	return p.WrapConnectionPgError(func() error {
		var dbPassword string
		query := p.pgPool.QueryRow(p.ctx, getUser, user.Login)
		if err := query.Scan(&user.ID, &dbPassword, &user.Balance, &user.Withdrawal); err != nil {
			logger.Debug(err.Error())
			var pgError *pgconn.PgError
			if errors.As(err, &pgError) {
				return err
			}
			return &customErr.UsersError{Code: 401, Err: "unknown login/password"}
		}
		logger.Debug(fmt.Sprintf("get user on db id:'%s' login:'%s", user.ID.String(), user.Login))
		return nil
	})
}

func (p *PostgreRepo) AddUser(user *models.DBUser) error {
	return p.WrapConnectionPgError(func() error {
		query := p.pgPool.QueryRow(p.ctx, insertUser, user.Login, user.Password)
		err := query.Scan(&user.ID)
		logger.Error(user.ID.String())
		if err != nil {
			var pgError *pgconn.PgError
			if errors.As(err, &pgError) && pgError.Code == pgerrcode.UniqueViolation {
				return &customErr.UsersError{Code: 409, Err: "login already exists"}
			}
			return err
		}
		return nil
	})
}

func (p *PostgreRepo) GetUserBalance(userID uuid.UUID) (*models.ViewUserBalance, error) {
	balance := models.ViewUserBalance{}
	err := p.WrapConnectionPgError(func() error {
		query := p.pgPool.QueryRow(p.ctx, getUserBalance, userID)
		if err := query.Scan(&balance.Current, &balance.Withdrawn); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &balance, err
}
