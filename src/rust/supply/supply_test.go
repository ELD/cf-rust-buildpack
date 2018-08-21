package supply_test

import (
	"io/ioutil"
	"path/filepath"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//go:generate mockgen -source=supply.go --destination=mocks_test.go --package=supply_test

var _ = Describe("Supply", func() {
	var (
		err      error
		buildDir string
	)
	BeforeEach(func() {
		buildDir, err = ioutil.TempDir("", "rust-buildpack.build.")
		Expect(err).To(BeNil())

		Describe("DetectCompilerVersion", func() {
			Context("Toolchain file exists", func() {
				It("takes from the toolchain file when it exists", func() {
					Expect(ioutil.WriteFile(filepath.Join(buildDir, "rustup-toolchain"), []byte("nightly-2018-08-18"), 0644)).To(Succeed())
				})
			})
		})
		It("example test", func() {
			Expect(false).To(Equal(false))
		})
		// TODO: Add tests here to check install dependency functions work
	})
})
