package errno

type Errno struct {
	HTTP    int
	Code    string
	Message string
}

func (err *Errno) Error() string {
	return err.Message
}

func (err *Errno) SetMessage(message string) *Errno {
	err.Message = message
	return err
}

func Decode(err error) (int, string, string) {
	if err == nil {
		return OK.HTTP, OK.Code, OK.Message
	}

	switch typed := err.(type) {
	case *Errno:
		return typed.HTTP, typed.Code, typed.Message
	default:
	}

	return InternalServerError.HTTP, InternalServerError.Code, err.Error()
}
