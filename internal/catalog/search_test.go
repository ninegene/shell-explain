package catalog

import (
	"strings"
	"testing"

	"github.com/ninegene/shell-explain/content"
)

func syntheticCatalog(t *testing.T) *Catalog {
	t.Helper()
	c := &Catalog{Entries: []Entry{
		{Metadata: Metadata{ID: "alpha", Title: "alpha"}, Sections: []Section{
			{Implementation: "common", Title: "Examples", Markdown: "Pack files into an archive.\n\n```sh\nalpha -c out.pk dir\n```\n\nList the Archive contents.\n\n```sh\nalpha -t out.pk\n```\n"},
			{Implementation: "gnu", Title: "GNU examples", Markdown: "Pack files into an archive.\n\n```sh\nalpha -c out.pk dir\n```\n"},
		}},
		{Metadata: Metadata{ID: "beta", Title: "beta"}, Sections: []Section{
			{Implementation: "common", Title: "Examples", Markdown: "Show the version.\n\n```sh\nbeta --archive-mode\n```\n\nCopy everything.\n\n```sh\nbeta cp a b\n```\n"},
			{Implementation: "common", Title: "Overview", Markdown: "Overview prose with archive.\n"},
		}},
	}}
	return c
}

func TestSearchExamples(t *testing.T) {
	c := syntheticCatalog(t)
	hits := c.SearchExamples("ARCHIVE", 0)
	if len(hits) != 3 {
		t.Fatalf("want 3 hits (dedupe across implementations, no overview), got %d: %+v", len(hits), hits)
	}
	if hits[0].Code != "alpha -c out.pk dir" || hits[0].Implementation != "common" {
		t.Errorf("first hit should be the deduplicated common example: %+v", hits[0])
	}
	if hits[2].Command != "beta" {
		t.Errorf("code-only match must rank last: %+v", hits)
	}
	if got := c.SearchExamples("archive contents", 0); len(got) != 1 || got[0].Code != "alpha -t out.pk" {
		t.Errorf("all words must match: %+v", got)
	}
	if got := c.SearchExamples("archive", 1); len(got) != 1 {
		t.Errorf("limit not applied: %d", len(got))
	}
	if got := c.SearchExamples("  ", 0); len(got) != 0 {
		t.Errorf("blank query must match nothing: %+v", got)
	}
	if got := c.SearchExamples("zzz", 0); got == nil || len(got) != 0 {
		t.Errorf("no match must be an empty, non-nil slice: %#v", got)
	}
}

func TestSearchExamplesEmbeddedContent(t *testing.T) {
	c, err := Load(content.Files)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, h := range c.SearchExamples("extract archive", 0) {
		if h.Command == "tar" {
			found = true
		}
		if strings.TrimSpace(h.Code) == "" {
			t.Errorf("hit without code: %+v", h)
		}
	}
	if !found {
		t.Error("\"extract archive\" must return tar examples")
	}
}
