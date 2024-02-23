package entities

type Result[T any] struct {
	Error  string
	Result T
}

func (result Result[T]) IsError() bool {
	return len(result.Error) != 0
}

func Success[T any](arg T) Result[T] {
	return Result[T]{
		Result: arg,
	}
}

func Error[T any](message string) Result[T] {
	return Result[T]{
		Error: message,
	}
}
