package tlsx_test

import (
	"os"
	"testing"

	"github.com/retrovibed/retroapi/internal/testx"
)

func TestMain(m *testing.M) {
	testx.Logging()
	os.Exit(m.Run())
}
