package services

import (
	"github.com/jmoiron/sqlx"
)

type Provider struct {
	Db *sqlx.DB
}

func New(db *sqlx.DB) *Provider {
	return &Provider{
		Db: db,
	}
}
