package application

import (
	"context"
	"encoding/gob"

	"log"

	"go.aoe.com/flamingo/framework/web"
)

type (
	ContextService struct{}

	StateContext struct {
		CurrentState State
		Data         Data
	}

	Data struct {
	}

	State interface {
		Process(ctx context.Context) error
	}
)

func init() {
	gob.Register(StateContext{})
}

func (cs *ContextService) GetCheckoutContext(ctx web.Context) *StateContext {
	if _, ok := ctx.Session().Values["cart.checkout.context"]; !ok {
		ctx.Session().Values["cart.checkout.context"] = StateContext{}
	}
	if contextInSession, ok := ctx.Session().Values["cart.checkout.context"]; ok {
		contextType, ok := contextInSession.(StateContext)
		if !ok {
			log.Printf("cart.checkout.ContextService Error - no valid Context saved in Session")
		}
		return &contextType
	}
	log.Printf("cart.checkout.ContextService Error - Context cannot be saved in Session")
	return nil
}
