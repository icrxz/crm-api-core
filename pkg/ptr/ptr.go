package ptr

// ToString returns the value of a *string pointer, or an empty string if nil.
func ToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
