# td

td provides is a simpler implementation of testing/T.TempDir.

Standard TempDir uses T.Name() to constaract the temp directory, while this implementation is just piggubacking off of os.MkTempDir. This strategy aims to avoid too long directory name failure, due to long test name like `t.Run("too long test name ...", ...)`.

## Usage

Use TD.TempDir instead of testing T.TempDir.

```go
package some_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/naoyafurudono/td"
)

func TestSomeThing(t *testing.T) {
  d := td.New(t)
  dir := d.TempDir()
  ...
}
```
