package examplecheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUVInitAllowsDifferentProjectLayouts(t *testing.T) {
	dir := t.TempDir()
	project := filepath.Join(dir, "demo")
	if err := os.Mkdir(project, 0o755); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"README.md", "pyproject.toml"} {
		if err := os.WriteFile(filepath.Join(project, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	const code = "uv init -q demo && ls demo"
	for _, output := range []string{"README.md\nmain.py\npyproject.toml\n", "README.md\npyproject.toml\nsrc\n"} {
		if err := Verify(code, dir, output); err != nil {
			t.Errorf("valid layout %q: %v", output, err)
		}
	}
	if err := Verify(code, dir, "README.md\nmain.py\n"); err == nil {
		t.Error("accepted output without pyproject.toml")
	}
	if err := os.Remove(filepath.Join(project, "pyproject.toml")); err != nil {
		t.Fatal(err)
	}
	if err := Verify(code, dir, "README.md\npyproject.toml\n"); err == nil {
		t.Error("accepted missing project file")
	}
}
