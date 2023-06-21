package dataencoder

import (
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"github.com/peter-mount/go-script/tools/dataencoder"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Include copies the include directory to the distribution
type Include struct {
	Encoder *dataencoder.Encoder `kernel:"inject"`
	Build   *dataencoder.Build   `kernel:"inject"`
	Source  *string              `kernel:"flag,include,install includes"`
}

func (s *Include) Start() error {
	s.Build.AddLibProvider(s.includeInclude)
	s.Build.AddLibProvider(s.includeDemo)

	if *s.Source != "" {
		return walk.NewPathWalker().
			Then(s.copy).
			Walk(*s.Source)
	}

	return nil
}

func (s *Include) includeInclude(dest string) (string, []string) {
	return filepath.Join(dest, "demo"), []string{"-include", "include"}
}

func (s *Include) includeDemo(dest string) (string, []string) {
	return filepath.Join(dest, "demo"), []string{"-include", "demo"}
}

func (s *Include) copy(path string, info os.FileInfo) error {
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
