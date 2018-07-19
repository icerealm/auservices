package main

import (
	"auservices/api"
	"auservices/domain"
	"testing"

	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

func TestAddCommandCategoryMessage(t *testing.T) {
	db, mock, err := sqlmock.New()

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	msgSequence := 10000
	apiCategory := &api.Category{
		Name: "TstCategory",
	}
	whoUpdate := "UnitTest"
	mock.ExpectBegin()
	//(message_seq, info, status, create_dt, rev_by)

	mock.ExpectExec("INSERT INTO category_message").WithArgs(msgSequence, apiCategory.Name)
	mock.ExpectCommit()
	domain.CreateCategory(db, int64(msgSequence), apiCategory, whoUpdate)
}
