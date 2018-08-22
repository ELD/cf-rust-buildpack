package CommandExt_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCommandExt(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CommandExt Suite")
}
