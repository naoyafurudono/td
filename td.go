package td

import (
	"fmt"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

type TD struct {
	t *testing.T

	parent string
	seq    int
	err    error
	mu     sync.Mutex
}

func New(t *testing.T) *TD {
	return &TD{t: t}
}

// Behaves as if T.TempDir of standard testing package, expect for this implmentaion does not use t.Name() as the directory name to be created.
func (td *TD) TempDir() string {
	// As the standard implementation, use a single parent directory for all
	// the temporary directories created by a test, each numbered sequentially.
	td.mu.Lock()
	var nonExistent bool
	if td.parent == "" {
		nonExistent = true
	} else {
		_, err := os.Stat(td.parent)
		nonExistent = os.IsNotExist(err)
		if err != nil && !nonExistent {
			td.t.Fatalf("TempDir: %v", err)
		}
	}

	if nonExistent {
		td.t.Helper()

		td.parent, td.err = os.MkdirTemp("", "")
		if td.err == nil {
			td.t.Cleanup(func() {
				if err := removeAll(td.parent); err != nil {
					td.t.Errorf("TempDir RemoveAll cleanup: %v", err)
				}
			})
		}
	}

	if td.err == nil {
		td.seq++
	}
	seq := td.seq
	td.mu.Unlock()

	if td.err != nil {
		td.t.Fatalf("TempDir: %v", td.err)
	}

	dir := fmt.Sprintf("%s%c%03d", td.parent, os.PathSeparator, seq)
	if err := os.Mkdir(dir, 0777); err != nil {
		td.t.Fatalf("TempDir: %v", err)
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
