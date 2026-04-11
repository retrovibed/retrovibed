package backoffx_test

import (
	"os"

	"github.com/retrovibed/retrovibed/internal/testx"

	"testing"
)

func TestMain(m *testing.M) {
	testx.Logging()
	os.Exit(m.Run())
}
