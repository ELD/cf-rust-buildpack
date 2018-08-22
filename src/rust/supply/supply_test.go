package supply_test

import (
	"bytes"
	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/libbuildpack/ansicleaner"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
	"rust/supply"
)

//go:generate mockgen -source=supply.go -destination=mocks_test.go -package=supply_test

var _ = Describe("Supply", func() {
	var (
		err            error
		buildDir       string
		buffer         *bytes.Buffer
		mockCtrl       *gomock.Controller
		mockManifest   *MockManifest
		mockInstaller  *MockInstaller
		mockStager     *MockStager
		mockCommand    *MockCommand
		mockCommandExt *MockCommandExt
		logger         *libbuildpack.Logger
		supplier       *supply.Supplier
	)

	BeforeEach(func() {
		buildDir, err = ioutil.TempDir("", "rust-buildpack.build.")
		Expect(err).To(BeNil())

		mockCtrl = gomock.NewController(GinkgoT())

		mockManifest = NewMockManifest(mockCtrl)
		mockInstaller = NewMockInstaller(mockCtrl)
		mockStager = NewMockStager(mockCtrl)
		mockCommand = NewMockCommand(mockCtrl)
		mockCommandExt = NewMockCommandExt(mockCtrl)

		buffer = new(bytes.Buffer)
		logger = libbuildpack.NewLogger(ansicleaner.New(buffer))

		supplier = &supply.Supplier{
			Manifest:   mockManifest,
			Installer:  mockInstaller,
			Stager:     mockStager,
			Command:    mockCommand,
			CommandExt: mockCommandExt,
			Log:        logger,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	JustBeforeEach(func() {
		Expect(supplier.Setup()).To(Succeed())
	})

	Describe("DetectCompilerVersion", func() {
		Context("Toolchain file exists", func() {
			BeforeEach(func() {
				Expect(ioutil.WriteFile(filepath.Join(buildDir, "rust-toolchain"), []byte("nightly-2018-08-18\n"), 0644)).To(Succeed())
				mockStager.EXPECT().BuildDir().AnyTimes().Return(buildDir)
			})

			It("takes from the toolchain file when it exists", func() {
				err := supplier.DetectCompilerVersion()
				Expect(supplier.ToolchainVersion).To(Equal(" --default-toolchain nightly-2018-08-18"))
				Expect(err).To(BeNil())
			})
		})
		Context("Toolchain file doesn't exist", func() {
			BeforeEach(func() {
				mockStager.EXPECT().BuildDir().AnyTimes().Return(buildDir)
			})

			It("returns an empty string", func() {
				err := supplier.DetectCompilerVersion()
				Expect(supplier.ToolchainVersion).To(Equal(""))
				Expect(err).To(BeNil())
			})
		})
	})

	Describe("InstallCompiler", func() {
		Context("Uncached", func() {
			Context("custom toolchain version", func() {
				BeforeEach(func() {
					mockStager.EXPECT().BuildDir().Return(buildDir).AnyTimes()
				})

				It("installs the compiler with the default toolchain matching the toolchain file version", func() {
					// Arrange
					mockCommandExt.EXPECT().
						ExecuteWithPipe(
							mockStager.BuildDir(),
							os.Stdout,
							os.Stdout,
							"curl https://sh.rustup.rs -sSf",
							"sh -s -- -y --default-toolchain nightly-2018-08-18").
						Return(nil).
						Times(1)
					supplier.ToolchainVersion = " --default-toolchain nightly-2018-08-18"

					// Act
					err := supplier.InstallCompiler()

					// Assert
					Expect(err).To(BeNil())
				})
			})

			Context("Default toolchain version", func() {
				BeforeEach(func() {
					mockStager.EXPECT().BuildDir().Return(buildDir).AnyTimes()
				})

				It("installs the compiler defaulting to the latest stable version of the compiler", func() {
					// Arrange
					mockCommandExt.EXPECT().
						ExecuteWithPipe(
							mockStager.BuildDir(),
							os.Stdout,
							os.Stdout,
							"curl https://sh.rustup.rs -sSf",
							"sh -s -- -y").
						Return(nil).
						Times(1)

					// Act
					err := supplier.InstallCompiler()

					// Assert
					Expect(err).To(BeNil())
				})
			})
		})

		// TODO: IMPLEMENT ME
		Context("Cached", func() {})
	})

	Describe("Compile code", func() {

	})
})
