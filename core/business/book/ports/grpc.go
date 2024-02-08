package ports

import (
	// "context"
	// "encoding/json"
	// "errors"
	// "fmt"
	// "io/ioutil"
	// "net/http"
	// "strings"

	// "main.go/common"
	accv1 "main.go/proto/core/customer/v1"

	"main.go/business/customer/app"
	// "main.go/business/customer/app/command"
	// "main.go/business/customer/domain"
	// "main.go/services/crypto"
	// "main.go/services/logger"

	// "google.golang.org/api/oauth2/v2"

	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	// "google.golang.org/grpc/codes"
)

type GRPCService struct {
	*accv1.UnimplementedCustomerServiceServer
	app        app.CustomerApp
	commandBus *cqrs.CommandBus
}

func NewGRPCService(app app.CustomerApp, commandBus *cqrs.CommandBus) (*GRPCService, error) {
	var err error
	if err != nil {
		return nil, err
	}
	return &GRPCService{
		app:        app,
		commandBus: commandBus,
	}, nil
}

// func (s *GRPCService) CreateCustomer(ctx context.Context, request *accv1.UserCustomerRequest) (*accv1.UserCustomerResponse, error) {
// 	usr := domain.User{
// 		User: &accv1.User{
// 			Name:     request.User.Name,
// 			Email:    request.User.Email,
// 			Password: crypto.HashPassword(ctx, []byte(request.User.Password)),
// 		},
// 	}
// 	user, err := s.app.Commands.CreateCustomer.Handle(ctx, usr)
// 	if err != nil {
// 		logger.Infof(ctx, "Error creating user customer {%v}", user)
// 		return nil, fmt.Errorf("error on create user. %s", err)
// 	}
// 	return &accv1.UserCustomerResponse{
// 		User: &accv1.User{
// 			Name:  user.Name,
// 			Email: user.Email,
// 		},
// 	}, nil
// }
