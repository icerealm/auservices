package domain

import (
	"auservices/api"
	"auservices/utilities"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	_ "github.com/lib/pq"
)

//GetConnection get connection
func GetConnection(cfg utilities.Configuration) (*sql.DB, error) {
	return sql.Open(cfg.DbDriver, cfg.DbURL)
}

//CreateCategory create category
func CreateCategory(db *sql.DB, msgSequence int64, apiCategory *api.Category, whoUpdate string) error {
	if len(whoUpdate) == 0 {
		return errors.New("whoUpdate is required parameter and cannot be empty string")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	bjson, err := json.Marshal(apiCategory)
	if err != nil {
		return err
	}

	msg := &CommandMessage{
		messageSeq: msgSequence,
		info:       string(bjson),
		status:     MessagePending,
		msgType:    CategoryTypeMsg,
		revBy:      whoUpdate,
		createDt:   time.Now(),
	}
	if err = msg.Save(tx); err != nil {
		tx.Rollback()
		return err
	}

	category := &Category{
		categoryNm:   apiCategory.Name,
		categoryDesc: apiCategory.Description,
		categoryType: api.CategoryType_name[int32(apiCategory.Type)],
		userID:       apiCategory.User.Userid,
		revDt:        time.Now(),
		revBy:        whoUpdate,
	}
	if err = category.Save(tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
