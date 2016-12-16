package hydrocarbon

import (
	"testing"

	"github.com/surullabs/lint"
)

func TestLint(t *testing.T) {
	// Run default linters
	err := lint.Default.Check("./...")

	// Ignore lint errors from auto-generated files
	err = lint.Skip(err, lint.RegexpMatch(`internal/`, `_string\.go`, `\.pb\.go`))

	if err != nil {
		t.Fatalf("lint failures: %v\n", err)
	}
}
