package book


import (
	"main.go/business/book/app"
	"main.go/business/book/ports"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"main.go/services/logger"
	"google.golang.org/grpc"
)




func NewApplicationFromSettings() (app.BookApp, error) {
	return app.BookApp{
		CommandHandlers: CommandHandlers,
		EventHandlers:   EventHandlers,
	}, nil
}


func CommandHandlers(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler {
	return []cqrs.CommandHandler{
		app.NewBookRoomHandler(eventBus),
		app.NewOrderBeerHandler(eventBus),
	}
}

func EventHandlers(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler {
	return []cqrs.EventHandler{
		app.NewOrderBeerOnRoomBooked(commandBus),
		app.NewBookingsFinancialReport(),
	}
}

func RegisterGRPCServer(app app.BookApp, cb *cqrs.CommandBus, registrars ...grpc.ServiceRegistrar) {
	_ , err := ports.NewGRPCService(app, cb)
	if err != nil {
		logger.Panicf(nil, "Error NewGRPCService %v", err)
	}
}

