package app


import (
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"main.go/business/book/app/command"
)

type CommandHandlerFunction func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler
type EventHandlerFunction func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler

type BookApp struct {
	CommandHandlers CommandHandlerFunction
	EventHandlers   EventHandlerFunction
}


func NewBookRoomHandler(eventBus *cqrs.EventBus) cqrs.CommandHandler {
	return command.NewBookRoomHandler(eventBus)
}


func NewOrderBeerHandler(eventBus *cqrs.EventBus) cqrs.CommandHandler {
	return command.NewOrderBeerHandler(eventBus)
}


func NewOrderBeerOnRoomBooked(commandBus *cqrs.CommandBus) cqrs.EventHandler {
	return command.NewOrderBeerOnRoomBooked(commandBus)
}

func NewBookingsFinancialReport() cqrs.EventHandler {
	return command.NewBookingsFinancialReport()
}

