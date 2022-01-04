package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table words...")
		_, err := db.Exec(`CREATE TABLE words(
			id BIGSERIAL PRIMARY KEY,
      word TEXT NOT NULL UNIQUE,
      sentence TEXT,
      translate TEXT,
      note TEXT,
      page_title TEXT,
      page_url TEXT,
      favicon_url TEXT,
      extra JSON,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table words...")
		_, err := db.Exec(`DROP TABLE words`)
		return err
	})
}
