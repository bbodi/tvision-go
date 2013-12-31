package tvision_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTvision(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Tvision Suite")
}
