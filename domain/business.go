package domain

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	//CategoryTypeMsg = CATEGORY_MESSAGE
	CategoryTypeMsg = "CATEGORY_MESSAGE"
	//ItemLineTypeMsg = ITEM_MESSAGE
	ItemLineTypeMsg = "ITEM_MESSAGE"
	//MessagePending = PENDING
	MessagePending = "PENDING"
)

//CommandMessage represent message from message server
type CommandMessage struct {
	id         int64
	messageSeq int64
	info       string
	status     string
	msgType    string
	revBy      string
	createDt   time.Time
}

//Category represent category
type Category struct {
	id           int64
	categoryNm   string
	categoryDesc string
	categoryType string
	userID       string
	revDt        time.Time
	revBy        string
}

//Save save msg
func (mp *CommandMessage) Save(tx *sql.Tx) error {
	prepareQuery := `INSERT INTO %s(message_seq, info, status, create_dt, rev_by)
	VALUES($1, $2, $3, CURRENT_TIMESTAMP, $4) RETURNING id`
	var query string
	if mp.msgType == CategoryTypeMsg {
		query = fmt.Sprintf(prepareQuery, "category_message")
	} else {
		query = fmt.Sprintf(prepareQuery, "item_message")
	}
	return tx.QueryRow(query,
		mp.messageSeq,
		mp.info,
		mp.status,
		mp.revBy).Scan(&mp.id)
}

//Save persist Category to Category table
func (cat *Category) Save(tx *sql.Tx) error {
	query := `INSERT INTO category(category_nm, category_desc, category_type, user_id, rev_dt, rev_by) 
	VALUES($1, $2, $3, $4, CURRENT_TIMESTAMP, $5) RETURNING id`
	return tx.QueryRow(query,
		cat.categoryNm,
		cat.categoryDesc,
		cat.categoryType,
		cat.userID,
		cat.revBy).Scan(&cat.id)
}
