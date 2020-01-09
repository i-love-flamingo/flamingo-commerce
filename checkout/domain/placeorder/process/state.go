package process

type (
	//State interface
	State interface {
		SetProcess(ctx *Process)
		Run() (*RollbackReference, error)
		IsFinal() bool
	}
)
