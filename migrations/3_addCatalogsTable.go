package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table catalogs...")
		_, err := db.Exec(`CREATE TABLE catalogs(
			id SERIAL PRIMARY KEY,
			title TEXT NOT NULL,
			parent_id INT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
      user_id INT REFERENCES users ON DELETE CASCADE
		)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table catalogs...")
		_, err := db.Exec(`DROP TABLE catalogs`)
		return err
	})
}
