package storage

// import (
// 	"context"
// 	"fmt"
// 	"iFall/internal/config"

// 	"github.com/jackc/pgx/v5/pgxpool"
// )

// type Storage struct {
// 	Pool *pgxpool.Pool
// }

// func MustConnect(cfg config.StorageConfig) *Storage {
// 	ctx, cancel := context.WithTimeout(context.Background(), cfg.ConnectTimeout)
// 	defer cancel()
// 	connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&timezone=%s",
// 		cfg.User,
// 		cfg.Password,
// 		cfg.Host,
// 		cfg.Port,
// 		cfg.Database,
// 		cfg.SSLMode,
// 		cfg.Timezone,
// 	)
// 	pcfg, err := pgxpool.ParseConfig(connString)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to parse pool config: %w", err))
// 	}
// 	pcfg.MaxConns = cfg.AmountOfConns
// 	pool, err := pgxpool.NewWithConfig(ctx, pcfg)
// 	if err != nil {
// 		panic(fmt.Errorf("failed to create postgres pool: %w", err))
// 	}
// 	ctxPing, cancel := context.WithTimeout(context.Background(), cfg.PingTimeout)
// 	defer cancel()
// 	if err := pool.Ping(ctxPing); err != nil {
// 		panic(fmt.Errorf("failed to ping postgres: %w", err))
// 	}
// 	storage := &Storage{
// 		Pool: pool,
// 	}
// 	return storage
// }

// func (s *Storage) Close() {
// 	s.Pool.Close()
// }
