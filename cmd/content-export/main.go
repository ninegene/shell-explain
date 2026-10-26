// content-export validates authored files and writes the exact CLI model for Astro.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/ninegene/shell-explain/content"
	"github.com/ninegene/shell-explain/internal/catalog"
)

func main() {
	output := flag.String("output", "", "Output file (stdout by default)")
	source := flag.String("source", "", "Content directory to validate instead of embedded files")
	flag.Parse()
	var files fs.FS = content.Files
	if *source != "" {
		files = os.DirFS(*source)
	}
	c, err := catalog.Load(files)
	if err == nil {
		var data []byte
		data, err = json.MarshalIndent(c, "", "  ")
		if err == nil {
			data = append(data, '\n')
			if *output == "" {
				_, err = os.Stdout.Write(data)
			} else {
				err = os.MkdirAll(filepath.Dir(*output), 0755)
				if err == nil {
					err = writeAtomic(*output, data)
				}
			}
		}
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Readers (including Astro's watcher) must see a complete catalog on each save.
func writeAtomic(path string, data []byte) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".catalog-*")
	if err != nil {
		return err
	}
	defer os.Remove(f.Name())
	if _, err = f.Write(data); err != nil {
		f.Close()
		return err
	}
	if err = f.Close(); err != nil {
		return err
	}
	return os.Rename(f.Name(), path)
}
