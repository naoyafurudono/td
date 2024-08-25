//go:build !windows

package tmp

func isWindowsRetryable(err error) bool {
	return false
}