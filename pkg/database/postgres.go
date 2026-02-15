package database

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	pool *pgxpool.Pool
	once sync.Once
)

// InitDB inicializa a conexão com o banco de dados PostgreSQL
func InitDB() (*pgxpool.Pool, error) {
	var err error
	once.Do(func() {
		dsn := os.Getenv("DATABASE_URL")
		if dsn == "" {
			err = fmt.Errorf("DATABASE_URL não configurada")
			return
		}

		config, parseErr := pgxpool.ParseConfig(dsn)
		if parseErr != nil {
			err = fmt.Errorf("erro ao fazer parse da config do banco: %v", parseErr)
			return
		}

		pool, err = pgxpool.NewWithConfig(context.Background(), config)
		if err != nil {
			err = fmt.Errorf("erro ao conectar ao banco de dados: %v", err)
			return
		}

		if pingErr := pool.Ping(context.Background()); pingErr != nil {
			err = fmt.Errorf("erro ao pingar o banco de dados: %v", pingErr)
			return
		}
	})

	return pool, err
}

// GetPool retorna a instância do pool de conexões
func GetPool() *pgxpool.Pool {
	return pool
}

// CloseDB fecha a conexão com o banco de dados
func CloseDB() {
	if pool != nil {
		pool.Close()
	}
}
