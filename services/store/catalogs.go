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
	Count      int       `pg:"-" json:"count"`
}

type CatalogCount struct {
	CatalogId int `json:"catalogId"`
	Count     int `json:"count"`
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

	err := db.Model(&catalog).
		Where("user_id = ?", userId).
		Select()

	if err != nil {
		log.Error().Err(err).Msg("Error fetching user's catalogs")
		return nil, dbError(err)
	}

	return catalog, nil
}

func FetchUserCatalogAndCount(user *User) ([]CatalogCount, error) {
	var res []CatalogCount

	uwords := UWords{
		UserId: user.ID,
	}

	err := db.Model(&uwords).
		Column("catalog_id").
		ColumnExpr("count(*) AS count").
		Group("catalog_id", "user_id").
		Having("user_id=?", user.ID).
		Select(&res)

	if err != nil {
		return nil, err
	}
	return res, nil
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
