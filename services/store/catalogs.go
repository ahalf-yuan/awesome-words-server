package store

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

type Catalog struct {
	ID         int
	ParentId   int
	Title      string
	CreatedAt  time.Time
	ModifiedAt time.Time
	UserID     int `json:"-"`
}

func AddCatalogNode(user *User, catalog *Catalog) error {
	catalog.UserID = user.ID
	fmt.Printf("on ----:%d \n", user.ID)

	_, err := db.Model(catalog).Returning("*").Insert()
	if err != nil {
		fmt.Printf("on ----:%d, %d, %s \n", catalog.ID, catalog.ParentId, catalog.Title)
		log.Error().Err(err).Msg("Error inserting new catalog node")
	}
	return dbError(err)
}
