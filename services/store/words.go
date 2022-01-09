package store

import (
	"fmt"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog/log"
)

type Words struct {
	ID         int       `json:"id"`
	Word       string    `json:"word" from:"word" binding:"required"`
	Sentence   string    `json:"sentence"`
	Translate  string    `json:"translate"`
	Note       string    `json:"note"`
	PageTitle  string    `json:"pageTitle"`
	PageUrl    string    `json:"pageURL"`
	FaviconUrl string    `json:"faviconURL"`
	Extra      []string  `json:"extra"`
	CreatedAt  time.Time `json:"createTime"`
	ModifiedAt time.Time `json:"modifiedTime"`
	UserId     int       `pg:"-" json:"userId"`
	WordId     int       `pg:"-" json:"wordId"`
	CatalogId  int       `pg:"-" json:"catalogId"`
}

type UWords struct {
	ID        int `json:"id"`
	WordId    int `json:"wordId"`
	UserId    int `json:"userId"`
	CatalogId int `json:"catalogId"`
}

// 先查找 words 表，如果words表已存在，不需要写入，返回 word_id
// 如果words表已存在，需要写入，返回 word_id
// 写数据到关系表
func AddUserWord(user *User, word *Words) error {
	// insert word entity with doing nothing

	var _, err = db.Model(word).
		Column("id").
		Where("word = ?word").
		OnConflict("DO NOTHING"). // OnConflict is optional
		Returning("id").
		SelectOrInsert()

	if err != nil {
		log.Error().Err(err).Msg("Error selecting new word")
		return err
	}

	uwords := UWords{
		UserId:    user.ID,
		WordId:    word.ID,
		CatalogId: word.CatalogId,
	}

	// created - false 已存在该记录
	var _, errCreate = db.Model(&uwords).
		Where("user_id = ?", user.ID).
		Where("word_id = ?", word.ID).
		Where("catalog_id = ?", word.CatalogId). // 一个分类下可以插入多条
		OnConflict("DO NOTHING").
		Returning("id").
		SelectOrInsert()

	// fmt.Println("crested =>", created)

	if errCreate != nil {
		log.Error().Err(err).Msg("Error inserting new word")
		return err
	}

	return dbError(err)
}

func FetchUserWords(user *User) ([]Words, error) {
	userId := user.ID
	var words []Words

	err := db.Model(&words).
		ColumnExpr("words.*").
		ColumnExpr("a.user_id AS user_id, a.word_id AS word_id").
		Join("JOIN u_words AS a ON a.word_id = words.id").
		Where("user_id=?", userId).
		Select()

	if err != nil {
		log.Error().Err(err).Msg("Error fetching user's catalogs")
		return nil, dbError(err)
	}
	return words, nil
}

func DeleteWord(id int) error {
	uword := new(UWords)
	uword.ID = id

	_, err := db.Model(uword).WherePK().Delete()
	if err != nil {
		log.Error().Err(err).Msg("Error updating catalog")
	}
	return dbError(err)
}

// 结构体与数组一样，都是值传递，比如当把数组或结构体作为实参传给函数的形参时，会复制一个副本，所以为了提高性能，一般不会把数组直接传递给函数，而是使用切片(引用类型)代替，而把结构体传给函数时，可以使用指针结构体
// across wordId & userId to get the recordId in db
func FetchWordById(wordId int, userId int) (*UWords, error) {
	uword := &UWords{
		WordId: wordId,
		UserId: userId,
	}

	err := db.Model(uword).
		Where("word_id=?", wordId).
		Where("user_id=?", userId).
		Column("id").
		Select()

	if err == pg.ErrNoRows {
		fmt.Println("No Rows!!")
		log.Error().Err(err).Msg("No Rows")
		return nil, nil
	}

	if err != nil {
		log.Error().Err(err).Msg("Error fetching word in u_words")
		return nil, err
	}

	return uword, dbError(err)
}
