package au_test

import (
	"auservices/api"
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc"
)

func prepareConn() *grpc.ClientConn {
	dialPort := fmt.Sprintf(":%d", 7777)
	conn, err := grpc.Dial(dialPort, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Did not connect: %s", err)
	}
	return conn
}

func TestPing(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()

	c := api.NewPingClient(conn)

	msg := &api.PingMessage{Greeting: "hello by client"}
	resp, err := c.SayHello(context.Background(), msg)
	if err != nil {
		t.Errorf("Error when calling SayHello: %s", err)
	}
	equals(t, "beating...", resp.Greeting)
}

func TestFindCategories(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()
	c := api.NewCategoryServicesClient(conn)
	q := &api.CategoryQuery{Query: "test"}
	_, err := c.FindCatetories(context.Background(), q)
	ok(t, err)
	// t.Logf("FindCatetories resp:%v", resp)
}

var categoryName = "test insert 4"
var categoryDescription = "test"
var userId = "us1343"

func TestAddCetogory(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()
	c := api.NewCategoryServicesClient(conn)
	resp, err := c.AddCategory(context.Background(),
		&api.Category{
			Name:        categoryName,
			Description: categoryDescription,
			Type:        1,
			User:        &api.User{Userid: userId},
		})
	if err != nil {
		t.Errorf("Error when calling AddCategory: %s", err)
	}
	equals(t, &api.MsgResponse{ResponseMsg: "Created"}, resp)
}

func TestGetAllCategoryTypeValues(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()
	c := api.NewCategoryServicesClient(conn)
	_, err := c.GetAllCategoryTypeValues(context.Background(), &api.Empty{})
	ok(t, err)
	// t.Logf("GetAllCategoryTypeValues resp:%v", resp)
}

func TestGetCategoryByName(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()
	c := api.NewCategoryServicesClient(conn)
	query := &api.CategoryQuery{
		Query: fmt.Sprintf("name = %s", categoryName),
		User:  &api.User{Userid: userId},
	}
	cat, err := c.GetCategoryByName(context.Background(), query)
	ok(t, err)
	t.Logf("data = %+v", cat)
}

func TestAddItemLine(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()
	c := api.NewItemLineServiceClient(conn)
	category := api.Category{
		Name:        "test insert 4",
		Description: "test",
		Type:        1,
		User:        &api.User{Userid: "us1343"},
	}
	resp, err := c.AddItemLine(context.Background(),
		&api.ItemLine{
			ItemLineNm:   "item line 1",
			ItemLineDesc: "item line desc",
			ItemLineDt:   time.Now().Unix(),
			ItemValue:    200.25,
			Category:     &category,
		})
	if err != nil {
		t.Errorf("Error when calling AddItemLine: %s", err)
	}
	equals(t, &api.MsgResponse{ResponseMsg: "Created"}, resp)
}