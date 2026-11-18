package catalog

import "testing"

func TestParseExamplesPairsDescriptionsWithCode(t *testing.T) {
	markdown := "Show the version; reviewed only, because the output changes with each release:\n\n```sh\nssh -V\n```\n\n" +
		"Log in to a host; reviewed only because it connects:\n\n```sh\nssh user@host\n```\n\n" +
		"Check a `value`\nacross two lines:\n\n```sh\nprintf a\nprintf b\n```\n"
	got := ParseExamples(markdown)
	want := []Example{
		{"Show the version", "ssh -V"},
		{"Log in to a host", "ssh user@host"},
		{"Check a value across two lines", "printf a\nprintf b"},
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
