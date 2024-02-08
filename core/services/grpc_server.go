package services

import (
	"google.golang.org/grpc"
	"go.uber.org/zap"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
)

type ServerInfo struct {
	Server *grpc.Server
	Router *message.Router
}


func NewServer(settings *core.Settings, logger *zap.Logger) (*ServerInfo, error) {

	var server *grpc.Server

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	cqrsFacade, err := cqrs.NewFacade(newFacadeConfig(
		settings,
		redisstream.JSONMarshaler{},
		*publisher,
		*subscriber,
		watermillLogger,
		router,
		func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler {
			commandHandlers := []cqrs.CommandHandler{}
			commandHandlers = append(commandHandlers, customerv1.CommandHandlers(commandBus, eventBus)...)
			return commandHandlers
		},
		func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler {
			eventHandlers := []cqrs.EventHandler{}
			eventHandlers = append(eventHandlers, customerv1.EventHandlers(commandBus, eventBus)...)
			return eventHandlers
		},
	))
	if err != nil {
		return nil, err
	}
	return &ServerInfo{
		Server: server,
		Router: router,
	}, nil
}



func newFacadeConfig(
	settings *core.Settings,
	cqrsMarshaller cqrs.CommandEventMarshaler,
	publisher message.Publisher,
	subscriber message.Subscriber,
	logger watermill.LoggerAdapter,
	router *message.Router,
	commandHandlers func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler,
	eventHandlers func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler,
) cqrs.FacadeConfig {
	return cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
			return []cqrs.CommandHandler{
				BookRoomHandler{eb},
				OrderBeerHandler{eb},
			}
		},
		CommandsPublisher: commandsPublisher,
		CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			// we can reuse subscriber, because all commands have separated topics
			return commandsSubscriber, nil
		},
		GenerateEventsTopic: func(eventName string) string {
			// because we are using PubSub RabbitMQ config, we can use one topic for all events
			return "events"
	
			// we can also use topic per event type
			// return eventName
		},
		EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
			return []cqrs.EventHandler{
				OrderBeerOnRoomBooked{cb},
				NewBookingsFinancialReport(),
			}
		},
		EventsPublisher: eventsPublisher,
		EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
			config := amqp.NewDurablePubSubConfig(
				amqpAddress,
				amqp.GenerateQueueNameTopicNameWithSuffix(handlerName),
			)
	
			return amqp.NewSubscriber(config, logger)
		},
		Router:                router,
		CommandEventMarshaler: cqrsMarshaler,
		Logger:                logger,
	}
}