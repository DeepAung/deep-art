package utils

func NewCopy[T any](ptr *T) *T {
	tmp := *ptr
	return &tmp
}
