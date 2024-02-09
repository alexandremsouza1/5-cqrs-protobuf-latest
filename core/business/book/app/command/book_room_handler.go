package command


import (
	"context"
	"math/rand"
	"log"
	"time"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"main.go/proto/core/events"
)

type BookRoomHandler struct {
	eventBus *cqrs.EventBus
}

func (b BookRoomHandler) HandlerName() string {
	return "BookRoomHandler"
}

// NewCommand returns type of command which this handle should handle. It must be a pointer.
func (b BookRoomHandler) NewCommand() interface{} {
	return &events.BookRoom{}
}

func (b BookRoomHandler) Handle(ctx context.Context, c interface{}) error {
	// c is always the type returned by `NewCommand`, so casting is always safe
	cmd := c.(*events.BookRoom)

	// some random price, in production you probably will calculate in wiser way
	price := (rand.Int63n(40) + 1) * 10

	log.Printf(
		"Booked %s for %s from %s to %s",
		cmd.RoomId,
		cmd.GuestName,
		time.Unix(cmd.StartDate.Seconds, int64(cmd.StartDate.Nanos)),
		time.Unix(cmd.EndDate.Seconds, int64(cmd.EndDate.Nanos)),
	)

	// RoomBooked will be handled by OrderBeerOnRoomBooked event handler,
	// in future RoomBooked may be handled by multiple event handler
	if err := b.eventBus.Publish(ctx, &events.RoomBooked{
		ReservationId: watermill.NewUUID(),
		RoomId:        cmd.RoomId,
		GuestName:     cmd.GuestName,
		Price:         price,
		StartDate:     cmd.StartDate,
		EndDate:       cmd.EndDate,
	}); err != nil {
		return err
	}

	return nil
}

func NewBookRoomHandler(eventBus *cqrs.EventBus) cqrs.CommandHandler {
	return BookRoomHandler{eventBus}
}