package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPool interface {
	Close()
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
	CopyFrom(ctx context.Context, tableName pgx.Identifier, columnNames []string, rowSrc pgx.CopyFromSource) (int64, error)
	Ping(ctx context.Context) error
}

type Config struct {
	MaxPoolSize  int
	ConnAttempts int
	ConnTimeout  time.Duration
}
type Postgres struct {
	Pool *pgxpool.Pool
}

func New(ctx context.Context, dsn string, cfg *Config) (*Postgres, error) {
	// Значения по умолчанию
	def := Config{
		MaxPoolSize:  10,
		ConnAttempts: 1,
		ConnTimeout:  2 * time.Second,
	}
	if cfg == nil {
		cfg = &def
	}
	if cfg.ConnAttempts <= 0 {
		cfg.ConnAttempts = 1
	}
	if cfg.MaxPoolSize <= 0 {
		cfg.MaxPoolSize = def.MaxPoolSize
	}
	if cfg.ConnTimeout <= 0 {
		cfg.ConnTimeout = def.ConnTimeout
	}

	poolCfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("parse pg config: %w", err)
	}
	poolCfg.MaxConns = int32(cfg.MaxPoolSize)

	var pool *pgxpool.Pool
	for i := 0; i < cfg.ConnAttempts; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		pool, err = pgxpool.NewWithConfig(ctx, poolCfg)
		if err == nil {
			if pingErr := pool.Ping(ctx); pingErr == nil {
				return &Postgres{Pool: pool}, nil
			}
			pool.Close()
			err = fmt.Errorf("ping failed: %w", err)
		}

		// если это не последний ретрай — подождём таймаут или отмену контекста
		if i < cfg.ConnAttempts-1 {
			select {
			case <-time.After(cfg.ConnTimeout):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return nil, fmt.Errorf("failed to connect to postgres after %d attempts: %w", cfg.ConnAttempts, err)
}

func (p *Postgres) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
