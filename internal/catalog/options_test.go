package catalog

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestOptionsText(t *testing.T) {
	sections := []Section{
		{Options: []string{"-n"}, Markdown: "Print each matching line with its line number. This helps you locate a result in the original file."},
		{Title: "Overview", Markdown: "Not an option."},
		{Options: []string{"-i", "--ignore-case"}, Markdown: "Ignore letter case when matching."},
	}
	want := "-n\nPrint each matching line with its line number. This helps you locate a result in the original file.\n\n-i, --ignore-case\nIgnore letter case when matching.\n"
	if got := OptionsText(sections, 0); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	for _, width := range []int{24, 72, 100, 120, 0} {
		limit := width
		if limit <= 0 || limit > DescriptionWidth {
			limit = DescriptionWidth
		}
		for _, line := range strings.Split(OptionsText(sections, width), "\n") {
			if utf8.RuneCountInString(line) > limit {
				t.Errorf("width %d exceeded: %q", limit, line)
			}
		}
	}
}

func TestExampleDescriptionWidth(t *testing.T) {
	command := "printf '%s\\n' 'a  b'"
	markdown := strings.Repeat("A description with several words. ", 6) + "\n\n```sh\n" + command + "\n```\n"
	for _, line := range strings.Split(ExampleText(markdown), "\n") {
		if strings.HasPrefix(line, "# ") && utf8.RuneCountInString(line) > DescriptionWidth {
			t.Errorf("comment exceeded width: %q", line)
		}
	}
	if !strings.Contains(ExampleText(markdown), "\n"+command+"\n") {
		t.Fatal("command text changed")
	}
}
