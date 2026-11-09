package examplecheck

import "testing"

func TestTryItBlocks(t *testing.T) {
	md := "Intro.\n\n```sh\necho hello\n```\n\nYou should see this:\n\n```text\nhello\n```\n\nNext:\n\n```sh\ndf -h .\n```\n\nAnother, with a diagram in between:\n\n```sh\nps\n```\n\n```diagram\n a\n```\n\n```text\nnot its output\n```\n"
	got := TryItBlocks(md)
	if len(got) != 3 {
		t.Fatalf("blocks: %+v", got)
	}
	if got[0].Code != "echo hello\n" || !got[0].HasExpect || got[0].Expect != "hello\n" {
		t.Errorf("first: %+v", got[0])
	}
	if got[1].Code != "df -h .\n" || got[1].HasExpect {
		t.Errorf("second must have no expected output: %+v", got[1])
	}
	if got[2].HasExpect {
		t.Errorf("a text block after a diagram is not this command's output: %+v", got[2])
	}
	if err := CheckTryIt(got[0], "hello\n"); err != nil {
		t.Error(err)
	}
	if err := CheckTryIt(got[0], "bye\n"); err == nil {
		t.Error("a different output was accepted")
	}
	if err := CheckTryIt(got[1], "anything at all\n"); err != nil {
		t.Error("a block without expected output must accept any output")
	}
}

func TestUnsafeTryIt(t *testing.T) {
	for _, safe := range []string{
		"ls", "df -h .", "du -sh .", "uname -s", "ps -o pid,ppid,comm", "echo $$", "sleep 1",
		"mkdir practice\ncd practice\necho hello > hello.txt\ncat hello.txt",
		"ls -l ~/Documents",
	} {
		if got := UnsafeTryIt(safe); got != "" {
			t.Errorf("%q flagged as unsafe: %s", safe, got)
		}
	}
	// The check works by word and is deliberately conservative: it also flags
	// a risky word that is only text, so a guide rewords instead of arguing.
	if UnsafeTryIt("echo kill the process") != "kill" {
		t.Error("a risky word inside text should be flagged")
	}
	for code, want := range map[string]string{
		"sudo ls":                  "sudo",
		"rm -rf practice":          "rm",
		"ls && rm hello.txt":       "rm",
		"kill 1234":                "kill",
		"dd if=/dev/zero of=x":     "dd",
		"mv a b":                   "mv",
		"curl https://example.org": "curl",
		"echo x > /etc/hosts":      "a redirect to a system path",
		"echo x >/dev/sda":         "a redirect to a system path",
		"shutdown -h now":          "shutdown",
	} {
		if got := UnsafeTryIt(code); got != want {
			t.Errorf("UnsafeTryIt(%q) = %q, want %q", code, got, want)
		}
	}
}
