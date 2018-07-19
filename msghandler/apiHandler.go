package msghandler

import (
	"auservices/api"
	"encoding/json"
	"log"

	context "golang.org/x/net/context"
)

//Server represents gRPC server.
type Server struct {
	MsgPublisher  *MessagePublisher
	MsgSubscriber *MessageSubscriber
}

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *api.PingMessage) (*api.PingMessage, error) {
	log.Printf("Receive message %s", in.Greeting)
	return &api.PingMessage{Greeting: "bar"}, nil
}

// FindCatetories values
func (s *Server) FindCatetories(ctx context.Context, in *api.CategoryQuery) (*api.CategoryList, error) {
	log.Printf("Receive query message %s", in.Query)
	categories := []*api.Category{
		&api.Category{
			Name:        "Test001",
			Description: "Desc Test001",
			Type:        api.CategoryType(api.CategoryType_value["EXPENSE"]),
			User:        &api.User{Userid: "empapay0013er"},
		},
		&api.Category{
			Name:        "Test002",
			Description: "Desc Test002",
			Type:        api.CategoryType(api.CategoryType_value["INCOME"]),
			User:        &api.User{Userid: "gmarer0014er"},
		},
		&api.Category{
			Name:        "Test003",
			Description: "Desc Test003",
			Type:        api.CategoryType(api.CategoryType_value["INCOME"]),
			User:        &api.User{Userid: "berarer0015er"},
		},
	}
	return &api.CategoryList{
		Categories: categories,
	}, nil
}

// AddCategory add new category to store
func (s *Server) AddCategory(ctx context.Context, in *api.Category) (*api.CategoryResponse, error) {
	log.Printf("insert category with %v", *in)
	b, err := json.Marshal(in)
	if err != nil {
		return &api.CategoryResponse{
			ResponseMsg: "FAILED",
		}, err
	}
	c := make(chan ConfirmationMessage)
	fn := func(uid string, err error) {
		if err != nil {
			resp := ConfirmationMessage{
				response: "ERROR",
				err:      err,
			}
			c <- resp
		} else {
			resp := ConfirmationMessage{
				response: uid,
				err:      nil,
			}
			c <- resp
		}
	}
	s.MsgPublisher.PublishEvent(kcategoryChannelID, string(b), fn)

	if ret := <-c; ret.err != nil {
		return &api.CategoryResponse{
			ResponseMsg: "Error",
		}, ret.err
	}
	return &api.CategoryResponse{
		ResponseMsg: "Created",
	}, err
}

//GetAllCategoryTypeValues all category type enum values
func (s *Server) GetAllCategoryTypeValues(ctx context.Context, in *api.Empty) (*api.CategortyTypeValues, error) {
	return &api.CategortyTypeValues{
		Types: api.CategoryType_value,
	}, nil
}
