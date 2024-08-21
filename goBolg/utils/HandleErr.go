// utils/error_utils.go
package utils

import "fmt"

// HandleErr is a utility function to handle errors and return the result or a formatted error.
func HandleErr[T any](result T, err error, errMsg string) (T, error) {
	if err != nil {
		var zero T
		return zero, fmt.Errorf("%s: %w", errMsg, err)
	}
	return result, nil
}
