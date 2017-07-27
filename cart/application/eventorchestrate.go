package application

type (
	// EventOrchestration type
	EventOrchestration struct {
		Cartservice *CartService `inject:""`
	}
)

/*
// Notify Implement Subscriber Interface
func (s *EventOrchestration) Notify(ev event.Event) {
	fmt.Printf("Event disoatched to CartService %s", ev)

	switch ev := ev.(type) {
	case domain.LoginSucessEvent:
		s.Cartservice.OnLogin(ev)
	}
}
*/
