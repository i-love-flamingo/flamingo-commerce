package application

import (
	"flamingo/core/cart/domain"
	"flamingo/framework/event"
	"fmt"
)

type (
	EventOrchestration struct {
		Cartservice *Cartservice `inject:""`
	}
)

//Implement Subscriber Interface
func (s *EventOrchestration) Notify(ev event.Event) {
	fmt.Printf("Event disoatched to Cartservice %s", ev)

	switch ev := ev.(type) {
	case domain.LoginSucessEvent:
		s.Cartservice.OnLogin(ev)
	}
}
