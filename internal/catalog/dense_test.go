package catalog

import (
	"strings"
	"testing"
)

func TestParseExamplesPairsDescriptionsWithCode(t *testing.T) {
	markdown := "Show the version; reviewed only, because the output changes with each release:\n\n```sh\nssh -V\n```\n\n" +
		"Log in to a host; reviewed only because it connects:\n\n```sh short=\"Log in\"\nssh user@host\n```\n\n" +
		"Check a `value`\nacross two lines:\n\n```sh\nprintf a\nprintf b\n```\n"
	got := ParseExamples(markdown)
	want := []Example{
		{Description: "Show the version", Code: "ssh -V"},
		{Description: "Log in to a host", Short: "Log in", Code: "ssh user@host"},
		{Description: "Check a value across two lines", Code: "printf a\nprintf b"},
	}
	if len(got) != len(want) {
		t.Fatalf("got %d examples, want %d: %+v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("example %d: got %+v, want %+v", i, got[i], want[i])
		}
	}
}

func TestCleanDescriptionKeepsOtherSentences(t *testing.T) {
	cases := map[string]string{
		"Run it; reviewed only:":                         "Run it",
		"Run it, reviewed only, because it connects:":    "Run it",
		"Run it. The batch forms differ:":                "Run it. The batch forms differ",
		"Run it; reviewed only. The batch forms differ:": "Run it. The batch forms differ",
		"Print the version":                              "Print the version",
	}
	for in, want := range cases {
		if got := CleanDescription(in); got != want {
			t.Errorf("%q: got %q, want %q", in, got, want)
		}
	}
}

func TestShortLabel(t *testing.T) {
	good := map[string]string{"": "", `short="Log in"`: "Log in"}
	for in, want := range good {
		if got, err := ShortLabel(in); err != nil || got != want {
			t.Errorf("%q: got %q, %v", in, got, err)
		}
	}
	for _, bad := range []string{`short=""`, `short=" x"`, `short="a` + "`" + `b"`, `title="x"`, `short="x" extra`, `short="` + strings.Repeat("x", MaxShortLabel+1) + `"`} {
		if _, err := ShortLabel(bad); err == nil {
			t.Errorf("%q should be rejected", bad)
		}
	}
}
