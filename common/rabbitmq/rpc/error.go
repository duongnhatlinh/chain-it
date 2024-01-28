package rpc

type Error struct {
	Message string
}

func (e *Error) NewError(message string) {
	e.Message = message
}

func (e Error) IsNil() bool {

	if e.Message == "" {
		return true
	}

	return false
}
