
package command


import (
	"context"
	"math/rand"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"main.go/proto/core/events"
)



type OrderBeerOnRoomBooked struct {
	commandBus *cqrs.CommandBus
}

func (o OrderBeerOnRoomBooked) HandlerName() string {
	// this name is passed to EventsSubscriberConstructor and used to generate queue name
	return "OrderBeerOnRoomBooked"
}

func (OrderBeerOnRoomBooked) NewEvent() interface{} {
	return &events.RoomBooked{}
}

func (o OrderBeerOnRoomBooked) Handle(ctx context.Context, e interface{}) error {
	event := e.(*events.RoomBooked)

	orderBeerCmd := &events.OrderBeer{
		RoomId: event.RoomId,
		Count:  rand.Int63n(10) + 1,
	}

	return o.commandBus.Send(ctx, orderBeerCmd)
}

func NewOrderBeerOnRoomBooked(commandBus *cqrs.CommandBus) cqrs.EventHandler {
	return OrderBeerOnRoomBooked{commandBus}
}