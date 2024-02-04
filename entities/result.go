package entities

type Result[T any] struct {
	Error  string
	Result T
}

func (result Result[T]) IsError() bool {
	return len(result.Error) != 0
}
