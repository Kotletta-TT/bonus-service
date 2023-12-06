package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Kotletta-TT/bonus-service/config"
	"github.com/Kotletta-TT/bonus-service/internal/logger"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (p *PostgreRepo) WrapConnectionPgError(pgFunc func() error) error {
	var err error
	attempt := 3
	timeoutRerun := 1
	for i := 0; i <= attempt; i++ {
		err = pgFunc()
		if err == nil {
			break
		}
		var pgError *pgconn.PgError
		if errors.As(err, &pgError) {
			if pgerrcode.IsConnectionException(pgError.Code) {
				logger.Error("Connect DB", "error:", err.Error(), "attempt", fmt.Sprintf("%d", i+1))
				time.Sleep(time.Duration(time.Second * time.Duration(timeoutRerun)))
				timeoutRerun += 2
				continue
			}
		}
		return err
	}
	return err
}

type PostgreRepo struct {
	cfg    *config.Config
	pgPool *pgxpool.Pool
	ctx    context.Context
}

func WrapErrDBNew(cnf *config.Config) *PostgreRepo {
	db, err := prNew(cnf)
	if err != nil {
		panic(err)
	}
	return db
}

func prNew(cnf *config.Config) (*PostgreRepo, error) {
	ctx := context.Background()
	pgConn, err := pgxpool.New(ctx, cnf.DSN)
	if err != nil {
		return nil, err
	}
	err = pgConn.Ping(ctx)
	if err != nil {
		return nil, err
	}
	p := &PostgreRepo{
		pgPool: pgConn,
		cfg:    cnf,
		ctx:    ctx,
	}
	p.initDB()
	return p, nil
}

func (p *PostgreRepo) initDB() {
	p.WrapConnectionPgError(func() error {
		_, err := p.pgPool.Exec(p.ctx, initUUIDextension)
		if err != nil {
			return err
		}
		_, err = p.pgPool.Exec(p.ctx, createTableUsers)
		if err != nil {
			return err
		}
		_, err = p.pgPool.Exec(p.ctx, createTableOrders)
		if err != nil {
			return err
		}
		_, err = p.pgPool.Exec(p.ctx, createTableWithdrawals)
		if err != nil {
			return err
		}
		return nil
	})
}

func (p *PostgreRepo) Close() error {
	logger.Info("Close DB connection")
	return nil
}
