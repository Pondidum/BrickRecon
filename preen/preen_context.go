package preen

type PreenContext struct {
	LinkToController ControllerLinker
}

func (pc *PreenContext) Error(err error) Error {
	return Error{ErrorMessage: err.Error()}
}

func (pc *PreenContext) ErrorS(err string) Error {
	return Error{ErrorMessage: err}
}
