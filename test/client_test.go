package main

import (
	"auservices/api"
	"context"
	"fmt"
	"log"
	"testing"

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
	t.Logf("resp:%v", resp.Greeting)
}

func TestFindCategories(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()
	c := api.NewCategoryServicesClient(conn)
	q := &api.CategoryQuery{Query: "test"}
	resp, err := c.FindCatetories(context.Background(), q)
	if err != nil {
		t.Errorf("Error when calling FindCategories: %s", err)
	}
	t.Logf("FindCatetories resp:%v", resp)
}

func TestAddCetogory(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()
	c := api.NewCategoryServicesClient(conn)
	resp, err := c.AddCategory(context.Background(),
		&api.Category{
			Name:        "test insertx 3",
			Description: "test",
			Type:        1})
	if err != nil {
		t.Errorf("Error when calling AddCategory: %s", err)
	}
	t.Logf("AddCategory resp:%v", resp)
}

func TestGetAllCategoryTypeValues(t *testing.T) {
	conn := prepareConn()
	defer conn.Close()
	c := api.NewCategoryServicesClient(conn)
	resp, err := c.GetAllCategoryTypeValues(context.Background(), &api.Empty{})

	if err != nil {
		t.Errorf("Error when calling GetAllCategoryTypeValues: %s", err)

	}
	t.Logf("GetAllCategoryTypeValues resp:%v", resp)

}
