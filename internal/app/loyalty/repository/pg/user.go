package pg

import (
	"context"

	"github.com/Kotletta-TT/bonus-service/internal/app/loyalty/entity"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
)

// AddUser adds a new user to the PostgreRepo.
//
// It takes a pointer to a User struct and a context as parameters.
// It returns an error indicating if the operation was successful or not.
func (pg *PostgreRepo) AddUser(user *entity.User, ctx context.Context) error {
	logger.Debug("Adding user", "login:", user.Login)
	return pg.WrapRetryPgConnErr(func() error {
		query := pg.pool.QueryRow(
			ctx,
			`INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`,
			user.Login,
			user.Password,
		)
		err := query.Scan(&user.ID)
		if err != nil {
			logger.Debug("Adding user failed", "error:", err.Error())
			return err
		}
		logger.Debug("Adding user success", "id:", user.ID)
		return nil
	})
}

// GetUserByLogin retrieves a user from the PostgreRepo using their login.
//
// Parameters:
// - login: the login of the user to retrieve.
// - ctx: the context for the database operation.
//
// Returns:
// - *entity.User: the retrieved user.
// - error: if there was an error retrieving the user.
func (pg *PostgreRepo) GetUserByLogin(login string, ctx context.Context) (*entity.User, error) {
	logger.Debug("Get user by login", "login:", login)
	user := &entity.User{}
	query := pg.pool.QueryRow(ctx,
		`SELECT
			id, password, current, withdrawn
		FROM
			users
		WHERE
			login = $1
		LIMIT
			1`,
		user.Login)
	err := pg.WrapRetryPgConnErr(func() error {
		return query.Scan(&user.ID, &user.Password, &user.Balance, &user.Withdrawal)
	})
	if err != nil {
		logger.Debug("Get user by login failed", "error:", err.Error())
		return nil, err
	}
	logger.Debug("Get user by login success", "id:", user.ID)
	return user, nil
}
