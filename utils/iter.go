package utils

func Reversed[T any](s []T) []T {
	n := len(s)
	out := make([]T, n)
	for i := 0; i < n; i++ {
		out[i] = s[n-1-i]
	}
	return out
}
