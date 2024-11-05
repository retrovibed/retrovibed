package errorsx_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestErrorsx(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Errorsx Suite")
}
