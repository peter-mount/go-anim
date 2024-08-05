package font

import (
	"archive/zip"
	"bytes"
	"fmt"
	"github.com/peter-mount/go-build/application"
	"github.com/peter-mount/go-build/core"
	"github.com/peter-mount/go-build/util/arch"
	"github.com/peter-mount/go-build/util/makefile/target"
	"github.com/peter-mount/go-build/util/meta"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	goScriptFonts   = "https://github.com/peter-mount/go-anim/archive/refs/heads/main.zip"
	goScriptFontDir = "lib/font"
)

// FontDownloader provides an extension that will automatically download
// the fonts from the main branch of the repository.
//
// It's not needed for the main go-anim build but is for projects
// that use it will need this if they use fonts
type FontDownloader struct {
	Encoder *core.Encoder `kernel:"inject"`
	Build   *core.Build   `kernel:"inject"`
}

func (s *FontDownloader) Start() error {
	s.Build.AddExtension(s.extension)
	return nil
}

func (s *FontDownloader) extension(arch arch.Arch, target target.Builder, meta *meta.Meta) {
	destDir := filepath.Join(arch.BaseDir(*s.Encoder.Dest), application.FileName(application.STATIC, "lib"))

	target.
		Target(destDir).
		MkDir(destDir).
		Echo("DLOAD", destDir).
		BuildTool("-copydir", "lib", "-d", destDir)

	// Now check we have a lib directory
	info, err := os.Stat(goScriptFontDir)
	if err == nil && !info.IsDir() {
		err = fmt.Errorf("%q is not a directory", goScriptFontDir)
	}
	if err != nil && os.IsNotExist(err) {
		err = s.downloadFonts()
	}
	if err != nil {
		fmt.Println(err)
		os.Exit(99)
	}
}

func (s *FontDownloader) downloadFonts() error {
	fmt.Println("Downloading fonts")
	resp, err := http.Get(goScriptFonts)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	zr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return err
	}

	for _, zipFile := range zr.File {
		if strings.HasPrefix(zipFile.Name, "go-anim-main/lib/font") && !strings.HasSuffix(zipFile.Name, "/") {
			err = s.Unzip(zipFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *FontDownloader) Unzip(zipFile *zip.File) error {
	dest := strings.Join(strings.Split(zipFile.Name, "/")[1:], "/")
	fmt.Printf("%-8s %s\n", "INSTALL", dest)

	r, err := zipFile.Open()
	if err != nil {
		return err
	}
	defer r.Close()

	err = os.MkdirAll(filepath.Dir(dest), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, r)
	return err
}
