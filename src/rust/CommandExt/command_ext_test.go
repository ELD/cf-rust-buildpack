package CommandExt_test

import (
	"bytes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"rust/CommandExt"
)

//go:generate mockgen -source=command_ext.go -destination=mocks_test.go -package=CommandExt_test

var _ = Describe("CommandExt", func() {
	var (
		buffer   *bytes.Buffer
		cmd      CommandExt.CommandExt
		command1 string
		command2 string
	)

	BeforeEach(func() {
		buffer = new(bytes.Buffer)
	})

	Context("ExecuteWithPipe", func() {
		BeforeEach(func() {
			command1 = `echo "test"`
			command2 = `grep "test"`
		})

		It("runs the command with the output in the right location", func() {
			err := cmd.ExecuteWithPipe("", buffer, buffer, command1, command2)

			Expect(err).To(BeNil())
			Expect(buffer.String()).To(ContainSubstring(`"test"`))
		})
	})
})
