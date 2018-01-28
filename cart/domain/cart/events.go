package cart

import (
	"context"

	"go.aoe.com/flamingo/framework/web"
)

type (
	OrderPlacedEvent struct {
		Cart           *Cart
		OrderId        string
		CurrentContext web.Context
	}

	//EventPublisher - technology free interface which is used in the Domain Layer to publish events that might be interesting for outside (Publish)
	EventPublisher interface {
		PublishOrderPlacedEvent(ctx context.Context, cart *Cart, orderId string)
	}
)
