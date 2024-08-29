package utils

func Asserts(condition bool, message string) {
	if !condition {
		panic(message)
	}
}
