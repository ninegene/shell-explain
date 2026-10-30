package catalog

import "testing"

func TestExampleText(t *testing.T) {
	markdown := "Find lines containing **ERROR** in a log:\n\n```sh\ngrep 'ERROR' app.log\n```\n\nShow matching line numbers:\n\n```sh\ngrep -n -i 'error' app.log\n```\n"
	want := "# Find lines containing ERROR in a log:\ngrep 'ERROR' app.log\n\n# Show matching line numbers:\ngrep -n -i 'error' app.log\n"
	if got := ExampleText(markdown); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestExampleTextPreservesCodeAndCommentsEveryProseLine(t *testing.T) {
	markdown := "Use `printf` &amp; keep spaces.\nHard break here.  \nThis is another description line.\n\n```sh\nprintf '%s\\n' 'a  b'\nif true; then\n  printf done\nfi\n```\n\nA final note with [reference](https://example.org/manual).\n"
	want := "# Use printf & keep spaces. Hard break here. This is another description line.\nprintf '%s\\n' 'a  b'\nif true; then\n  printf done\nfi\n\n# A final note with reference (https://example.org/manual).\n"
	if got := ExampleText(markdown); got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got := ExampleText("```sh\nprintf one\n```\n\n```sh\nprintf two\n```\n"); got != "printf one\n\nprintf two\n" {
		t.Fatalf("adjacent examples: %q", got)
	}
}
