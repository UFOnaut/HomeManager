package errors

type LoginError struct {
	Message string
}

func (error LoginError) Error() string {
	return error.Message
}
