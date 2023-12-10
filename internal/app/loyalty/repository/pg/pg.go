package pg

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

type PostgreRepo struct {
	pool   *pgxpool.Pool
	config *config.Config
}

func New(config *config.Config) (*PostgreRepo, error) {
	pgConfig, err := pgxpool.ParseConfig(config.DSN)
	if err != nil {
		return nil, err
	}
	pgConfig.MaxConns = int32(10)
	pgConfig.MinConns = int32(5)
	pgPool, err := pgxpool.NewWithConfig(context.Background(), pgConfig)
	if err != nil {
		return nil, err
	}
	return &PostgreRepo{pool: pgPool, config: config}, nil
}

func (pg *PostgreRepo) WrapRetryPgConnErr(pgFunc func() error) error {
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
