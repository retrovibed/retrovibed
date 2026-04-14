package httpx_test

import (
	"os"
	"testing"

	"github.com/retrovibed/retrovibed/shallows/internal/testx"
)

func TestMain(m *testing.M) {
	testx.Logging()
	os.Exit(m.Run())
}
