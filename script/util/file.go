package util

import (
	"github.com/peter-mount/go-anim/util/frames"
	"github.com/peter-mount/go-kernel/v2/util/walk"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// GetImageFiles returns a list of image files from a directory.
// The results will be images which are supported and who's names are timestamps.
// The result will be sorted in the correct order.
func (_ *Util) GetImageFiles(dir string) ([]string, error) {
	var files []string
	err := walk.NewPathWalker().
		Then(func(path string, _ os.FileInfo) error {
			files = append(files, path)
			return nil
		}).
		If(isImage).
		IsFile().
		Walk(dir)
	if err != nil {
		return nil, err
	}

	sort.SliceStable(files, func(i, j int) bool {
		return files[i] < files[j]
	})

	return files, nil
}

func isImage(_ string, info os.FileInfo) bool {
	ext := filepath.Ext(info.Name())
	return in(ext, ".png", ".jpg", ".jpeg", ".tif", ".tiff")
}

func in(s string, p ...string) bool {
	s = strings.ToLower(s)
	for _, e := range p {
		if s == strings.ToLower(e) {
			return true
		}
	}
	return false
}

func (_ *Util) Sequence(interval int, sourceFiles []string) *frames.FrameSet {
	return frames.Sequence(interval, sourceFiles)
}

func (_ *Util) SequenceIn(interval int, sourceFiles []string, loc *time.Location) *frames.FrameSet {
	return frames.SequenceIn(interval, sourceFiles, loc)
}
