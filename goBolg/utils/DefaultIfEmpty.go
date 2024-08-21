package utils

// DefaultIfEmpty checks if the input string is empty, and returns the default value if it is.
func DefaultIfEmpty(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
