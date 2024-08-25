package tmp

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

type A struct {
	t *testing.T

	parent string
	seq    int
	err    error
	mu     sync.Mutex
}

func New(t *testing.T) *A {
	return &A{t: t}
}

func (a *A) TempDir() string {
	// As the standard implementation, use a single parent directory for all
	// the temporary directories created by a test, each numbered sequentially.
	a.mu.Lock()
	var nonExistent bool
	if a.parent == "" {
		nonExistent = true
	} else {
		_, err := os.Stat(a.parent)
		nonExistent = os.IsNotExist(err)
		if err != nil && !nonExistent {
			a.t.Fatalf("TempDir: %v", err)
		}
	}

	if nonExistent {
		a.t.Helper()

		a.parent, a.err = os.MkdirTemp("", "")
		if a.err == nil {
			a.t.Cleanup(func() {
				if err := removeAll(a.parent); err != nil {
					a.t.Errorf("TempDir RemoveAll cleanup: %v", err)
				}
			})
		}
	}

	if a.err == nil {
		a.seq++
	}
	seq := a.seq
	a.mu.Unlock()

	if a.err != nil {
		a.t.Fatalf("TempDir: %v", a.err)
	}

	dir := fmt.Sprintf("%s%c%03d", a.parent, os.PathSeparator, seq)
	if err := os.Mkdir(dir, 0777); err != nil {
		a.t.Fatalf("TempDir: %v", err)
	}
	return dir
}

// removeAll is like os.RemoveAll, but retries Windows "Access is denied."
// errors up to an arbitrary timeout.
//
// Those errors have been known to occur spuriously on at least the
// windows-amd64-2012 builder (https://go.dev/issue/50051), and can only occur
// legitimately if the test leaves behind a temp file that either is still open
// or the test otherwise lacks permission to delete. In the case of legitimate
// failures, a failing test may take a bit longer to fail, but once the test is
// fixed the extra latency will go away.
func removeAll(path string) error {
	const arbitraryTimeout = 2 * time.Second
	var (
		start     time.Time
		nextSleep = 1 * time.Millisecond
	)
	for {
		err := os.RemoveAll(path)
		if !isWindowsRetryable(err) {
			return err
		}
		if start.IsZero() {
			start = time.Now()
		} else if d := time.Since(start) + nextSleep; d >= arbitraryTimeout {
			return err
		}
		time.Sleep(nextSleep)
		nextSleep += time.Duration(rand.Int63n(int64(nextSleep)))
	}
}
