package api

import (
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
func (s *Server) SayHello(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Receive message %s", in.Greeting)
	return &PingMessage{Greeting: "bar"}, nil
}

// FindCatetories values
func (s *Server) FindCatetories(ctx context.Context, in *CategoryQuery) (*CategoryList, error) {
	log.Printf("Receive query message %s", in.Query)
	categories := []*Category{
		&Category{
			Cid:         "001",
			Name:        "Test001",
			Description: "Desc Test001",
			Type:        CategoryType(CategoryType_value["EXPENSE"]),
		},
		&Category{
			Cid:         "002",
			Name:        "Test002",
			Description: "Desc Test002",
			Type:        CategoryType(CategoryType_value["INCOME"]),
		},
		&Category{
			Cid:         "003",
			Name:        "Test003",
			Description: "Desc Test003",
			Type:        CategoryType(CategoryType_value["INCOME"]),
		},
	}
	return &CategoryList{
		Categories: categories,
	}, nil
}

// AddCategory add new category to store
func (s *Server) AddCategory(ctx context.Context, in *Category) (*CategoryResponse, error) {
	log.Printf("insert category with %v", *in)
	b, err := json.Marshal(in)
	if err != nil {
		return &CategoryResponse{
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
		return &CategoryResponse{
			ResponseMsg: "Error",
		}, ret.err
	}
	return &CategoryResponse{
		ResponseMsg: "Created",
	}, err
}

//GetAllCategoryTypeValues all category type enum values
func (s *Server) GetAllCategoryTypeValues(ctx context.Context, in *Empty) (*CategortyTypeValues, error) {
	return &CategortyTypeValues{
		Types: CategoryType_value,
	}, nil
}
