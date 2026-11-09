package content

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Every Markdown file under en must be in the embedded filesystem, so the
// offline CLI and the website are built from the same authored collection.
func TestEveryContentFileIsEmbedded(t *testing.T) {
	count := 0
	err := filepath.WalkDir("en", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(path, ".md") {
			return err
		}
		count++
		disk, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		embedded, err := Files.ReadFile(filepath.ToSlash(path))
		if err != nil {
			t.Errorf("%s is not embedded: add its directory to the go:embed directive", path)
			return nil
		}
		if string(disk) != string(embedded) {
			t.Errorf("%s differs from the embedded copy; rebuild", path)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if count == 0 {
		t.Fatal("no content files found")
	}
}
