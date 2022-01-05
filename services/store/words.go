package store

import (
	"time"

	"github.com/rs/zerolog/log"
)

type Words struct {
	ID         int       `json:"id"`
	Word       string    `json:"word" binding:"required"`
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
}

type UWords struct {
	ID     int `json:"id"`
	WordId int `json:"word_id"`
	UserId int `json:"user_id"`
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
		log.Error().Err(err).Msg("Error inserting new word")
		return err
	}

	uwords := UWords{
		UserId: user.ID,
		WordId: word.ID,
	}

	// created - false 已存在该记录
	var _, errCreate = db.Model(&uwords).
		Where("user_id = ?", user.ID).
		Where("word_id = ?", word.ID).
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
