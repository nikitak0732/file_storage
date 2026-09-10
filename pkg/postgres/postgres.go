package postgres

import (
	"fmt"
	"test_astral/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Postgres struct {
	Postgres *sqlx.DB
}

func New(cfg *config.Config) (*Postgres, error) {

	ConnString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.DbName,
	)

	db, err := sqlx.Open("postgres", ConnString)

	if err != nil {

		return nil, err
	}

	if err := db.Ping(); err != nil {

		return nil, err
	}

	return &Postgres{db}, nil
}

func (p *Postgres) Close() error {
	if p.Postgres != nil {

		return p.Postgres.Close()
	}
	return nil
}
