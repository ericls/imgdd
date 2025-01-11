package utils

func SafeDeref[T any](p *T) T {
	if p == nil {
		var v T
		return v
	}
	return *p
}

func SafeDerefWithDefault[T any](p *T, defaultValue T) T {
	if p == nil {
		return defaultValue
	}
	return *p
}
