package utils

func Must(val any, err error) any {
	if err != nil {
		panic(err)
	}
	return val
}

func ReturnFirst[T any](first T, second any) T {
	return first
}
