package services

import (
	"google.golang.org/grpc"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"

	"main.go/business/book"
)

type ServerInfo struct {
	Server *grpc.Server
	Router *message.Router
}


func NewServer(amqpAddress string) (*cqrs.Facade, error) {

	var server *grpc.Server
	logger := watermill.NewStdLogger(false, false)

	bookv1 , err := book.NewApplicationFromSettings()
	if err != nil {
		return nil, err
	}

	cqrsMarshaler := cqrs.ProtobufMarshaler{}
	
	// You can use any Pub/Sub implementation from here: https://watermill.io/docs/pub-sub-implementations/
	// Detailed RabbitMQ implementation: https://watermill.io/docs/pub-sub-implementations/#rabbitmq-amqp
	// Commands will be send to queue, because they need to be consumed once.
	commandsAMQPConfig := amqp.NewDurableQueueConfig(amqpAddress)
	commandsPublisher, err := amqp.NewPublisher(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}
	commandsSubscriber, err := amqp.NewSubscriber(commandsAMQPConfig, logger)
	if err != nil {
		panic(err)
	}

	router, err := message.NewRouter(message.RouterConfig{}, logger)
	if err != nil {
		panic(err)
	}

	// Events will be published to PubSub configured Rabbit, because they may be consumed by multiple consumers.
	// (in that case BookingsFinancialReport and OrderBeerOnRoomBooked).
	eventsPublisher, err := amqp.NewPublisher(amqp.NewDurablePubSubConfig(amqpAddress, nil), logger)
	if err != nil {
		panic(err)
	}

	cqrsFacade, err := cqrs.NewFacade(newFacadeConfig(
		func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler {
				// command.BookRoomHandler{eventBus},
				// command.OrderBeerHandler{eventBus},

				commandHandlers := []cqrs.CommandHandler{}
			  commandHandlers = append(commandHandlers, bookv1.CommandHandlers(commandBus, eventBus)...)
				return commandHandlers
		},
		func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler {
				// command.OrderBeerOnRoomBooked{commandBus},
				// command.NewBookingsFinancialReport(),
				eventHandlers := []cqrs.EventHandler{}
				eventHandlers = append(eventHandlers, bookv1.EventHandlers(commandBus, eventBus)...)
				return eventHandlers
		},
		eventsPublisher,
		commandsPublisher,
		commandsSubscriber,
		router,
		cqrsMarshaler,
		logger,
		func(handlerName string) (message.Subscriber, error) {
				config := amqp.NewDurablePubSubConfig(
				amqpAddress,
				amqp.GenerateQueueNameTopicNameWithSuffix(handlerName),
			)

			return amqp.NewSubscriber(config, logger)
		},

	))
	if err != nil {
		return nil, err
	}

	book.RegisterGRPCServer(bookv1, cqrsFacade.CommandBus(), server)


	return cqrsFacade,nil
}



func newFacadeConfig(
	commandHandlers func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler,
	eventHandlers func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler,
	eventsPublisher message.Publisher,
	commandsPublisher message.Publisher,
	commandsSubscriber message.Subscriber,
	router *message.Router,
	cqrsMarshaler cqrs.CommandEventMarshaler,
	logger watermill.LoggerAdapter,
	eventsSubscriberConstructor func(handlerName string) (message.Subscriber, error),
) cqrs.FacadeConfig {
	return cqrs.FacadeConfig{
		GenerateCommandsTopic: func(commandName string) string {
			// we are using queue RabbitMQ config, so we need to have topic per command type
			return commandName
		},
		// CommandHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.CommandHandler {
		// 	return []cqrs.CommandHandler{
		// 		BookRoomHandler{eb},
		// 		OrderBeerHandler{eb},
		// 	}
		// },
		CommandHandlers: commandHandlers,
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
		// EventHandlers: func(cb *cqrs.CommandBus, eb *cqrs.EventBus) []cqrs.EventHandler {
		// 	return []cqrs.EventHandler{
		// 		OrderBeerOnRoomBooked{cb},
		// 		NewBookingsFinancialReport(),
		// 	}
		// },
		EventHandlers: eventHandlers,
		EventsPublisher: eventsPublisher,
		// EventsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
		// 	config := amqp.NewDurablePubSubConfig(
		// 		amqpAddress,
		// 		amqp.GenerateQueueNameTopicNameWithSuffix(handlerName),
		// 	)

		// 	return amqp.NewSubscriber(config, logger)
		// },
		EventsSubscriberConstructor: eventsSubscriberConstructor,
		Router:                router,
		CommandEventMarshaler: cqrsMarshaler,
		Logger:                logger,
	}
}