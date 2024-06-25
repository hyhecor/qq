package qq

func first[T any](aa ...T) (b T) {
	for _, v := range aa {
		return v
	}
	return
}

func String[T ~string](v *T) string {
	if v == nil {
		return ""
	}
	return (string)(*v)
}
