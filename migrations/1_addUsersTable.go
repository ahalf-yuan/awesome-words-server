package main

import (
	"fmt"

	"github.com/go-pg/migrations/v8"
)

func init() {
	// type: 'email','wechat',...
	// third_unique_acount: 第三方唯一用户id，可以是微信的openid，可以是QQ的openid，抑或苹果id
	// username - 统一邮箱登陆
	// nick_name - 邮箱前缀
	migrations.MustRegisterTx(func(db migrations.DB) error {
		fmt.Println("creating table users...")
		_, err := db.Exec(`CREATE TABLE users(
			id SERIAL PRIMARY KEY,
			username TEXT NOT NULL UNIQUE,
			hashed_password BYTEA NOT NULL,
			salt BYTEA NOT NULL,
      avatar TEXT,
      nick_name VARCHAR(255),
      type VARCHAR(64) NOT NULL,
      third_unique_acount TEXT,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
			modified_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)`)
		return err
	}, func(db migrations.DB) error {
		fmt.Println("dropping table users...")
		_, err := db.Exec(`DROP TABLE users`)
		return err
	})
}
