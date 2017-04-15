package domain

import "flamingo/core/event2"

type(
	LoginSucessEvent struct {
		event2.DefaultEvent
		UserId string
	}
)

//Factory for Event
func NewLoginSucessEvent(userId string) LoginSucessEvent {
	return LoginSucessEvent{
		event2.DefaultEvent{"login.sucess"},
		userId,
	}
}
