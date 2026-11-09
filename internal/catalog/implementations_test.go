package catalog

import (
	"strings"
	"testing"
)

func fixture(id, aliases, impls, extra string) []byte {
	var b strings.Builder
	b.WriteString("---\nschema_version: 1\nid: " + id + "\ntitle: " + id + "\nsummary: Fixture entry.\naliases: " + aliases + "\ntags: [test]\nimplementations: " + impls + "\nreferences:\n  - https://example.org/manual\n---\n")
	for _, s := range []string{"Overview", "Background", "Examples", "Common mistakes", "Portability"} {
		b.WriteString("\n### " + s + "\n\nText.\n")
	}
	_ = extra
	return []byte(strings.Replace(b.String(), "\n### Overview", "\n## Common\n\n### Overview", 1))
}

func TestRegistryIsConsistent(t *testing.T) {
	seenID, seenHeading := map[string]bool{}, map[string]bool{}
	for _, impl := range Implementations {
		if !names.MatchString(impl.ID) || impl.Heading == "" || impl.Label == "" || seenID[impl.ID] || seenHeading[impl.Heading] {
			t.Fatalf("bad registry entry %+v", impl)
		}
		seenID[impl.ID], seenHeading[impl.Heading] = true, true
	}
	if Implementations[0].ID != "common" {
		t.Fatal("common must be first")
	}
}

func TestSymbolicAliasAndNewImplementations(t *testing.T) {
	e, err := Parse(fixture("test", "['[']", "[gnu]", "") /* needs GNU view */)
	if err == nil {
		t.Fatalf("missing declared GNU view accepted: %+v", e)
	}
	data := append(fixture("test", "['[']", "[upstream]", ""), []byte("\n## Upstream\n\n### About upstream\n\nText.\n\n### Upstream examples\n\nText.\n")...)
	e, err = Parse(data)
	if err != nil {
		t.Fatal(err)
	}
	c := &Catalog{Entries: []Entry{*e}}
	if got, ok := c.Find("["); !ok || got.ID != "test" {
		t.Fatal("symbolic alias not found")
	}
	if _, err := e.Select("upstream", ""); err != nil {
		t.Fatal(err)
	}
	zsh := append(fixture("where", "[]", "[bash, zsh]", ""), []byte("\n## Bash\n\n### About Bash\n\nText.\n\n### Bash examples\n\nText.\n\n## Zsh\n\n### About Zsh\n\nText.\n\n### Zsh examples\n\nText.\n")...)
	z, err := Parse(zsh)
	if err != nil || !z.ShellOnly() || len(z.ShellImplementations()) != 2 {
		t.Fatal(err, z)
	}
	if e.ShellOnly() {
		t.Fatal("upstream entry is not shell-only")
	}
}

func TestRejectsUnsafeNamesAndMixedUpstream(t *testing.T) {
	for _, tt := range []struct{ id, aliases, impls string }{
		{"[", "[]", "[upstream]"},          // symbolic spellings are aliases, never IDs
		{"test", "['[[']", "[upstream]"},   // only registered symbolic spellings
		{"test", "['../x']", "[upstream]"}, // path characters
		{"test", "[]", "[upstream, gnu]"},  // upstream is a single-implementation view
		{"test", "[]", "[freebsd]"},        // unregistered
		{"test", "[]", "[common]"},         // common is implicit
		{"Test", "[]", "[upstream]"},       // IDs start with a lowercase letter
	} {
		if _, err := Parse(fixture(tt.id, tt.aliases, tt.impls, "")); err == nil {
			t.Errorf("accepted %+v", tt)
		}
	}
}
