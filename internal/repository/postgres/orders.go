package postgres

import (
	"github.com/Kotletta-TT/bonus-service/internal/errors"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/Kotletta-TT/bonus-service/internal/models"
	"github.com/google/uuid"
)

const setNewOrder = `
	INSERT INTO orders (user_id, number)
	VALUES ($1, $2)`

const getUserIDByOrderNumber = `
	SELECT 
		CASE 
			WHEN EXISTS (SELECT user_id FROM orders WHERE number = $1) 
				THEN (SELECT user_id FROM orders WHERE number = $1)
			ELSE NULL 
		END AS user_id_result`

const getUserOrders = `
	SELECT user_id, number, accrual, status, uploaded_at
	FROM orders
	WHERE user_id = $1
	ORDER BY uploaded_at
	DESC`

const getUnpocessedOrders = `
	SELECT user_id, number
	FROM orders 
	WHERE status 
	IN ('NEW', 'REGISTERED', 'PROCESSING')`

const updateOrder = `
	UPDATE orders 
	SET accrual = $1, status = $2
	WHERE number = $3`

func (p *PostgreRepo) GetUserOrders(userID uuid.UUID) ([]*models.DBOrders, error) {
	orders := make([]*models.DBOrders, 0, 10)
	err := p.WrapConnectionPgError(func() error {
		query, err := p.pgPool.Query(p.ctx, getUserOrders, userID)
		if err != nil {
			return err
		}
		for query.Next() {
			dbOrder := models.DBOrders{}
			err = query.Scan(&dbOrder.UserID, &dbOrder.Number, &dbOrder.Accrual, &dbOrder.Status, &dbOrder.UploadedAt)
			if err != nil {
				return err
			}
			orders = append(orders, &dbOrder)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (p *PostgreRepo) SetUserOrder(order *models.DBOrders) error {
	return p.WrapConnectionPgError(func() error {
		tx, err := p.pgPool.Begin(p.ctx)
		if err != nil {
			return err
		}
		orderDB := models.DBOrders{}
		existOrderQuery := tx.QueryRow(p.ctx, getUserIDByOrderNumber, order.Number)
		err = existOrderQuery.Scan(&orderDB.UserID)
		if err != nil {
			logger.Error(err.Error())
			return errors.InternalServerErr()
		}
		switch orderDB.UserID {
		case uuid.Nil:
			_, err = tx.Exec(p.ctx, setNewOrder, order.UserID, order.Number)
			if err != nil {
				logger.Error(err.Error())
				tx.Rollback(p.ctx)
				return errors.InternalServerErr()
			}
			tx.Commit(p.ctx)
			return nil
		case order.UserID:
			tx.Rollback(p.ctx)
			return errors.UploadOrderEarlier()
		default:
			return errors.UploadAnotherUser()
		}
	})
}

func (p *PostgreRepo) GetUnprocessedOrders() ([]*models.DBOrders, error) {
	ordersNum := make([]*models.DBOrders, 0, 10)
	err := p.WrapConnectionPgError(func() error {
		rows, err := p.pgPool.Query(p.ctx, getUnpocessedOrders)
		if err != nil {
			return err
		}
		for rows.Next() {
			order := models.DBOrders{}
			err := rows.Scan(&order.UserID, &order.Number)
			if err != nil {
				return err
			}
			ordersNum = append(ordersNum, &order)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ordersNum, nil
}

func (p *PostgreRepo) UpdateOrderAndBalance(order *models.DBOrders) error {
	return p.WrapConnectionPgError(func() error {
		trx, err := p.pgPool.Begin(p.ctx)
		if err != nil {
			return err
		}
		_, err = trx.Exec(p.ctx, updateOrder, order.Accrual, order.Status, order.Number)
		if err != nil {
			trx.Rollback(p.ctx)
			return err
		}
		row := trx.QueryRow(p.ctx, getUserBalance, order.UserID)
		usrBalance := models.ViewUserBalance{}
		err = row.Scan(&usrBalance.Current, &usrBalance.Withdrawn)
		if err != nil {
			trx.Rollback(p.ctx)
			return err
		}
		if order.Status == "PROCESSED" {
			usrBalance.Current += order.Accrual
			_, err = trx.Exec(p.ctx, updateUserBalance, usrBalance.Current, usrBalance.Withdrawn, order.UserID)
			if err != nil {
				trx.Rollback(p.ctx)
				return err
			}
		}
		trx.Commit(p.ctx)
		return nil
	})
}
