package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table u_words...")
		_, err := db.Exec(`CREATE TABLE u_words(
			id BIGSERIAL PRIMARY KEY,
      user_id INT REFERENCES users ON DELETE CASCADE,
      word_id INT REFERENCES words ON DELETE CASCADE,
      catalog_id INT DEFAULT -1 REFERENCES catalogs ON DELETE CASCADE,
      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table u_words...")
		_, err := db.Exec(`DROP TABLE u_words`)
		return err
	})
}
