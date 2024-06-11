package utils

func Map[T, Q any](arr []T, fun func(T) Q) []Q {
	newarr := make([]Q, len(arr))
	for i := range arr {
		newarr[i] = fun(arr[i])
	}

	return newarr
}

func Filter[T any](arr []T, fun func(T) bool) []T {
	var newarr []T
	for i := range arr {
		if fun(arr[i]) {
			newarr = append(newarr, arr[i])
		}
	}

	return newarr
}
