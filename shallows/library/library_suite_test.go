package library_test

import (
	"os"
	"testing"

	"github.com/retrovibed/retrovibed/internal/testx"
)

func TestMain(m *testing.M) {
	testx.Logging()
	os.Exit(m.Run())
}
