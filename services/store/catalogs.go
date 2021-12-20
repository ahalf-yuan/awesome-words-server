package store

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Catalog struct {
	ID         int       `json:"id"`
	ParentId   int       `json:"parentId"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"createTime"`
	ModifiedAt time.Time `json:"modifiedTime"`
	UserID     int       `json:"-"`
}

func AddCatalogNode(user *User, catalog *Catalog) error {
	catalog.UserID = user.ID

	_, err := db.Model(catalog).Returning("*").Insert()
	if err != nil {
		log.Error().Err(err).Msg("Error inserting new catalog node")
	}
	return dbError(err)
}

func FetchUserCatalogs(user *User) ([]Catalog, error) {
	userId := user.ID
	var catalog []Catalog

	err := db.Model(catalog).
		Where("user_id = ?", userId).
		Select()

	if err != nil {
		log.Error().Err(err).Msg("Error fetching user's catalogs")
		return nil, dbError(err)
	}
	return catalog, nil
}

func FetchCatalog(id int) (*Catalog, error) {
	catalog := new(Catalog)
	catalog.ID = id
	err := db.Model(catalog).WherePK().Select()
	if err != nil {
		log.Error().Err(err).Msg("Error fetching catalog")
		return nil, dbError(err)
	}
	return catalog, nil
}

func UpdateCatalog(catalog *Catalog) error {
	_, err := db.Model(catalog).WherePK().UpdateNotZero()
	if err != nil {
		log.Error().Err(err).Msg("Error updating catalog")
	}
	return dbError(err)
}

func DeleteCatalog(catalog *Catalog) error {
	_, err := db.Model(catalog).WherePK().Delete()
	if err != nil {
		log.Error().Err(err).Msg("Error deleting catalog")
	}
	return dbError(err)
}
