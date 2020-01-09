package process

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
func (f *Factory) New(ctx *Context) Process {
	if ctx == nil {
		ctx = &Context{
			State: f.startState,
		}
	}
	return Process{context: *ctx}
}

//Run triggers run on current state
func (p *Process) Run() {

}

//Context to get current process context
func (p *Process) Context() Context {
	return p.context
}

//UpdateState updates
func (p *Process) UpdateState(s State) {
	p.context.State = s
}
