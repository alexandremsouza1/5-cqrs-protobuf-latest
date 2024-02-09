package command

import (
	"context"
	"math/rand"
	"log"
	"github.com/pkg/errors"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"main.go/proto/core/events"
)

type OrderBeerHandler struct {
	eventBus *cqrs.EventBus
}

func (o OrderBeerHandler) HandlerName() string {
	return "OrderBeerHandler"
}

func (o OrderBeerHandler) NewCommand() interface{} {
	return &events.OrderBeer{}
}

func (o OrderBeerHandler) Handle(ctx context.Context, c interface{}) error {
	cmd := c.(*events.OrderBeer)

	if rand.Int63n(10) == 0 {
		// sometimes there is no beer left, command will be retried
		return errors.Errorf("no beer left for room %s, please try later", cmd.RoomId)
	}

	if err := o.eventBus.Publish(ctx, &events.BeerOrdered{
		RoomId: cmd.RoomId,
		Count:  cmd.Count,
	}); err != nil {
		return err
	}

	log.Printf("%d beers ordered to room %s", cmd.Count, cmd.RoomId)
	return nil
}

func NewOrderBeerHandler(eventBus *cqrs.EventBus) cqrs.CommandHandler {
	return OrderBeerHandler{eventBus}
}