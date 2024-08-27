//go:build !windows

package td

func isWindowsRetryable(err error) bool {
	return false
}