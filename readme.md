# tmp

tmp provides more stable version of testing/T.TempDir by taking simple implementation.

Standard TempDir uses T.Name() to constaract the temp directory, while this implementation is just piggubacking of os.MkTempDir. This strategy aims to avoid too long directory name failure, due to long test name.

## Usage

Use A.TempDir instead of testing T.TempDir.

```go
package some_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/naoyafurudono/tmp"
)

func TestSomeThing(t *testing.T) {
  a := tmp.New(t)
  dir := a.TempDir()
  ...
}
```
