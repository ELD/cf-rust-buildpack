package supply

import (
	"io"

	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cloudfoundry/libbuildpack"
)

type Stager interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/stager.go
	BuildDir() string
	DepDir() string
	DepsIdx() string
	DepsDir() string
}

type Manifest interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/manifest.go
	AllDependencyVersions(string) []string
	DefaultVersion(string) (libbuildpack.Dependency, error)
}

type Installer interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/installer.go
	InstallDependency(libbuildpack.Dependency, string) error
	InstallOnlyVersion(string, string) error
}

type Command interface {
	//TODO: See more options at https://github.com/cloudfoundry/libbuildpack/blob/master/command.go
	Execute(string, io.Writer, io.Writer, string, ...string) error
	Output(dir string, program string, args ...string) (string, error)
}

type Supplier struct {
	Manifest              Manifest
	Installer             Installer
	Stager                Stager
	Command               Command
	Log                   *libbuildpack.Logger
	appHasCargoTomlExists bool
	appHasCargoLockExists bool
}

func (s *Supplier) Run() error {
	s.Log.BeginStep("Supplying Rust")

	if err := s.Setup(); err != nil {
		return fmt.Errorf("Error during setup: %v", err)
	}

	version, err := s.DetectCompilerVersion()
	if err != nil {
		return fmt.Errorf("Unable ")
	}

	s.Command.Execute(
		s.Stager.BuildDir(),
		os.Stdout,
		os.Stdout,
		"curl",
		"https://sh.rustup.rs -sSf",
		"|",
		"sh",
		"--",
		"-y",
		version)

	return nil
}

func (s *Supplier) Setup() error {
	// Detect Cargo.toml and Cargo.lock
	if exists, err := libbuildpack.FileExists(filepath.Join(s.Stager.BuildDir(), "Cargo.toml")); err != nil {
		return fmt.Errorf("Unable to determine if Cargo.toml exists: %v", err)
	} else {
		s.appHasCargoTomlExists = exists
	}

	if exists, err := libbuildpack.FileExists(filepath.Join(s.Stager.BuildDir(), "Cargo.lock")); err != nil {
		return fmt.Errorf("Unable to determine if Cargo.lock exists: %v", err)
	} else {
		s.appHasCargoLockExists = exists
	}

	return nil
}

func (s *Supplier) DetectCompilerVersion() (string, error) {
	exists, _ := libbuildpack.FileExists(filepath.Join(s.Stager.BuildDir(), "rustup-toolchain"))

	toolchainVersion := ""
	if exists {
		bytes, err := ioutil.ReadFile("rustup-toolchain")

		if err != nil {
			return "", fmt.Errorf("Unable to read from 'rustup-toolchain' file: %v", err)
		}

		toolchainVersion = "--default-toolchain " + string(bytes)
	}

	return toolchainVersion, nil
}
