package build

import (
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/tools/build"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Lib copies the include directory to the distribution
type Lib struct {
	Encoder *core.Encoder `kernel:"inject"`
	Build   *build.Build  `kernel:"inject"`
	Source  *string       `kernel:"flag,lib,install lib"`
}

func (s *Lib) Start() error {
	s.Build.AddLibProvider(s.includeLib)

	if *s.Source != "" {
		return walk.NewPathWalker().
			Then(s.copy).
			Walk(*s.Source)

	}

	return nil
}

func (s *Lib) includeLib(dest string) (string, []string) {
	return filepath.Join(dest, "lib"), []string{"-lib", "lib"}
}

func (s *Lib) copy(path string, info os.FileInfo) error {
	// Ignore the source base directory
	if path == *s.Source {
		return nil
	}

	// dest is the source minus the source base directory name
	dstName := filepath.Join(*s.Encoder.Dest, strings.TrimPrefix(path, *s.Source+"/"))
	if info.IsDir() {
		return os.MkdirAll(dstName, info.Mode())
	}

	srcFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.OpenFile(dstName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, info.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
