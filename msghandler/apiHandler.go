package msghandler

import (
	"auservices/api"
	"auservices/domain"
	"auservices/utilities"
	"encoding/json"
	"log"

	context "golang.org/x/net/context"
)

//Server represents gRPC server.
type Server struct {
	MsgPublisher *MessagePublisher
}

// SayHello generates response to a Ping request
func (s *Server) SayHello(ctx context.Context, in *api.PingMessage) (*api.PingMessage, error) {
	log.Printf("Receive message %s", in.Greeting)
	return &api.PingMessage{Greeting: "beating..."}, nil
}

// FindCatetories values
func (s *Server) FindCatetories(ctx context.Context, in *api.CategoryQuery) (*api.CategoryList, error) {
	log.Printf("Receive query message %v", in)
	db, err := domain.GetConnection(utilities.GetConfiguration())
	if err != nil {
		log.Println("FindCatetories, error while connecting to database:", err)
		return nil, err
	}
	defer db.Close()
	categories, err := domain.GetCategories(db, in)
	if err != nil {
		log.Println("FindCatetories, error while retrieving category data:", err)
		return nil, err
	}
	return &api.CategoryList{
		Categories: categories,
	}, nil
}

// AddCategory add new category to store
func (s *Server) AddCategory(ctx context.Context, in *api.Category) (*api.MsgResponse, error) {
	log.Printf("insert category with %v", *in)
	b, err := json.Marshal(in)
	if err != nil {
		return &api.MsgResponse{
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
		return &api.MsgResponse{
			ResponseMsg: "Error",
		}, ret.err
	}
	return &api.MsgResponse{
		ResponseMsg: "Created",
	}, err
}

//GetAllCategoryTypeValues all category type enum values
func (s *Server) GetAllCategoryTypeValues(ctx context.Context, in *api.Empty) (*api.CategortyTypeValues, error) {
	return &api.CategortyTypeValues{
		Types: api.CategoryType_value,
	}, nil
}

//GetCategoryByName get category info by cateogry name
func (s *Server) GetCategoryByName(ctx context.Context, in *api.CategoryQuery) (*api.Category, error) {
	db, err := domain.GetConnection(utilities.GetConfiguration())
	if err != nil {
		log.Fatalf("GetCategoryByName, error while connecting to database: %v", err)
	}
	defer db.Close()
	return domain.GetCategoryByName(db, in)
}

// AddItemLine add new itemline to store
func (s *Server) AddItemLine(ctx context.Context, in *api.ItemLine) (*api.MsgResponse, error) {
	log.Printf("insert itemLine with %v", *in)
	b, err := json.Marshal(in)
	if err != nil {
		return &api.MsgResponse{
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
	s.MsgPublisher.PublishEvent(kitemLineChannelID, string(b), fn)

	if ret := <-c; ret.err != nil {
		return &api.MsgResponse{
			ResponseMsg: "Error",
		}, ret.err
	}
	return &api.MsgResponse{
		ResponseMsg: "Created",
	}, err
}
