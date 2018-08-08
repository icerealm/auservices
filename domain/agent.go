package domain

/**
// This file contains business logic that transaction control is here.
*/
import (
	"auservices/api"
	"auservices/utilities"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//GetConnection get connection
func GetConnection(cfg utilities.Configuration) (*sql.DB, error) {
	return sql.Open(cfg.DbDriver, cfg.DbURL)
}

func getQueryValue(key string, query string) string {
	s := strings.Split(query, "=")
	if len(s) > 0 && strings.Trim(s[0], " ") == key {
		return strings.Trim(s[1], " ")
	}
	return ""
}

//GetCategories get categories by userId
func GetCategories(db *sql.DB, query *api.CategoryQuery) ([]*api.Category, error) {
	user := query.User
	if user == nil {
		return []*api.Category{}, nil
	}
	repo := Repository{dbConn: db}
	categories, err := repo.FindAllCategories(user.Userid)
	if err != nil {
		return nil, err
	}
	var categoryList []*api.Category
	for _, category := range categories {
		categoryList = append(categoryList, &api.Category{
			Name:        category.categoryNm,
			Description: category.categoryDesc,
			Type:        api.CategoryType(api.CategoryType_value[category.categoryType]),
			User:        &api.User{Userid: category.userID},
		})
	}
	return categoryList, nil
}

//GetCategoryByName get category by name
func GetCategoryByName(db *sql.DB, query *api.CategoryQuery) (*api.Category, error) {
	name := getQueryValue("name", query.Query)
	user := query.User
	repo := Repository{dbConn: db}
	category, err := repo.FindCategoryByName(name, user.Userid)
	if err != nil {
		return nil, err
	}
	return &api.Category{
		Name:        category.categoryNm,
		Description: category.categoryDesc,
		Type:        api.CategoryType(api.CategoryType_value[category.categoryType]),
		User:        &api.User{Userid: category.userID},
	}, nil
}

//CreateCategory create category
func CreateCategory(db *sql.DB, msgSequence uint64, apiCategory *api.Category, whoUpdate string) error {
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

//CreateItemLine create itemline.
func CreateItemLine(db *sql.DB, msgSequence uint64, apiItemLine *api.ItemLine, whoUpdate string) error {
	if len(whoUpdate) == 0 {
		return errors.New("whoUpdate is required parameter and cannot be empty string")
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	bjson, err := json.Marshal(apiItemLine)
	if err != nil {
		return err
	}

	msg := &CommandMessage{
		messageSeq: msgSequence,
		info:       string(bjson),
		status:     MessagePending,
		msgType:    ItemLineTypeMsg,
		revBy:      whoUpdate,
		createDt:   time.Now(),
	}
	if err = msg.Save(tx); err != nil {
		tx.Rollback()
		return err
	}
	repository := &Repository{
		dbConn: db,
	}
	category, err := repository.FindCategoryByName(apiItemLine.Category.Name, apiItemLine.Category.User.Userid)
	if err != nil {
		tx.Rollback()
		return err
	}
	itemLine := &ItemLine{
		itemLineNm:    apiItemLine.ItemLineNm,
		itemLineDesc:  apiItemLine.ItemLineDesc,
		itemLineDt:    time.Unix(apiItemLine.ItemLineDt, 0),
		itemLineValue: apiItemLine.ItemValue,
		categoryID:    category.id,
		userID:        apiItemLine.Category.User.Userid,
		revBy:         whoUpdate,
	}
	if err = itemLine.Save(tx); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
