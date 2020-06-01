package preen

type Redirect struct {
	URL string
}

type Error struct {
	ErrorMessage string
}

func ErrorModel(err error) Error {
	return Error{ErrorMessage: err.Error()}
}
