package preen

type Error struct {
	ErrorMessage string
}

func ErrorModel(err error) Error {
	return Error{ErrorMessage: err.Error()}
}

func ErrorModelS(err string) Error {
	return Error{ErrorMessage: err}
}
