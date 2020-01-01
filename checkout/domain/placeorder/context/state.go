package context

type (
	Rollback func() error

	State interface {
		SetContext(ctx *Context)
		Run() (Rollback, error)
		IsFinal() bool
	}
)
