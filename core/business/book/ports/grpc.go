package ports

import (

	"main.go/business/book/app"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type GRPCService struct {
	app        app.BookApp
	commandBus *cqrs.CommandBus
}

func NewGRPCService(app app.BookApp, commandBus *cqrs.CommandBus) (*GRPCService, error) {
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
