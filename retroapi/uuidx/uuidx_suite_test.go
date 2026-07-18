package uuidx_test

import (
	"os"

	"github.com/retrovibed/retrovibed/retroapi/testx"

	"testing"
)

func TestMain(m *testing.M) {
	testx.Logging()
	os.Exit(m.Run())
}
