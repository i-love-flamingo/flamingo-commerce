package process

import "errors"

type (

	//Process representing a place order process and has a current context with infos about result and current state
	Process struct {
		context Context
	}

	//Factory use to get Process instance
	Factory struct {
		startState State
	}

	//RollbackReference a reference that can be used to trigger a rollback
	RollbackReference struct {
	}
)

//New returns new process - optional with given Context
func (f *Factory) Inject(dep *struct {
	StartState State `inject:"startState"`
}) {
	if dep != nil {
		f.startState = dep.StartState
	}
}

//New returns new process - optional with given Context
func (f *Factory) New(ctx *Context) (*Process, error) {
	if ctx == nil {
		if f.startState == nil {
			return nil, errors.New("No start state given")
		}
		ctx = &Context{
			State: f.startState,
		}
	}
	return &Process{context: *ctx}, nil
}

//Run triggers run on current state
func (p *Process) Run() {
	_ = p.context.State.Run(p)
}

//Context to get current process context
func (p *Process) Context() Context {
	return p.context
}

//UpdateState updates
func (p *Process) UpdateState(s State) {
	p.context.State = s
}
