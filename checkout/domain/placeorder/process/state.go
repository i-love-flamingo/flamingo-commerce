package process

type (
	//State interface
	State interface {
		Run(*Process) *RollbackReference
		IsFinal() bool
		Name() string
	}
)
