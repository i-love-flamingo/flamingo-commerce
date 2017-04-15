package application

import (
	"flamingo/core/event2"
	"fmt"
	"flamingo/core/cart/domain"
)



type (
	EventOrchestration struct {
		*event2.Subscriber `inject:""`
		Cartservice *Cartservice `inject:""`
	}
)


//Implement Subscribtions Interface
func (s *EventOrchestration) AddSubscriptions() {
	s.EventDispatcher.Subscribe("login.sucess", s)
}



//Implement Listerner Interface
func (s *EventOrchestration) OnEvent(event event2.Event) {
	fmt.Printf("Event disoatched to Cartservice %s",event)
	switch key := event.GetEventKey(); key {
	case "login.sucess":
		s.Cartservice.OnLogin(event.(domain.LoginSucessEvent))
	default:
		// Todo log waring?
		fmt.Printf("Unkonwn event for this package %s", key)
	}
}

