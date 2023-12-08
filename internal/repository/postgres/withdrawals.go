package postgres

import (
	e "errors"

	"github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

const getUserCurrentWithdraw = `
	SELECT COALESCE(SUM(sum), 0)
	FROM withdraws
	WHERE user_id = $1`

const setUserWithdraw = `
	INSERT INTO withdrawals (user_id, order_id, sum)
	VALUES ($1, $2, $3)`

const getUserWithdraws = `
	SELECT order_id, sum, processed_at
	FROM withdrawals
	WHERE user_id = $1
	ORDER BY processed_at
	DESC`

func (p *PostgreRepo) GetUserWithdraws(userID uuid.UUID) ([]*models.ViewWithdraw, error) {
	userWithdraws := make([]*models.ViewWithdraw, 0, 10)
	err := p.WrapConnectionPgError(func() error {
		query, err := p.pgPool.Query(p.ctx, getUserWithdraws, userID)
		if err != nil {
			return err
		}
		for query.Next() {
			withdraw := models.ViewWithdraw{}
			err = query.Scan(&withdraw.OrderID, &withdraw.Sum, &withdraw.ProcessedAt)
			if err != nil {
				return err
			}
			userWithdraws = append(userWithdraws, &withdraw)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return userWithdraws, nil
}

func (p *PostgreRepo) RequestUserWithdraw(withdraw *models.DBWithdraw) error {
	return p.WrapConnectionPgError(func() error {
		tx, err := p.pgPool.Begin(p.ctx)
		if err != nil {
			return err
		}
		// Получение текущих данных из таблицы Users
		usr := models.DBUser{ID: withdraw.UserID}
		queryOrder := tx.QueryRow(p.ctx, getUserBalance, withdraw.UserID)
		err = queryOrder.Scan(&usr.Balance, &usr.Withdrawal)
		if err != nil {
			return err
		}
		if withdraw.Sum > usr.Balance {
			return errors.EnoughtBalance()
		}
		// Попытка вставки данных в таблицу withdrawals
		_, err = tx.Exec(p.ctx, setUserWithdraw, withdraw.UserID, withdraw.OrderID, withdraw.Sum)
		if err != nil {
			tx.Rollback(p.ctx)
			var pgError *pgconn.PgError
			if e.As(err, &pgError) {
				if pgError.Code == pgerrcode.ForeignKeyViolation {
					return errors.IncorrectOrderNumber()
				}
			}
			// Обработка 2-х типов ошибок (violates foreigin key 23503, connection error)
			return err
		}
		// Расчеты для обновления данных
		usr.Balance -= withdraw.Sum
		usr.Withdrawal += withdraw.Sum
		// Обновление данных в таблице users
		_, err = tx.Exec(p.ctx, updateUserBalance, usr.Balance, usr.Withdrawal, usr.ID)
		if err != nil {
			tx.Rollback(p.ctx)
			return err
		}
		tx.Commit(p.ctx)
		return nil
	})
}
