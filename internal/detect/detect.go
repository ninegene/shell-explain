// Package detect identifies executable implementations without evaluating shell input.
package detect

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"
)

// probed lists the commands whose native executables may run `--version` to
// identify an implementation. Add a command only after checking that GNU,
// BSD/Apple and BusyBox builds reject or answer the flag without side effects
// and without waiting for input; do not probe commands outside this list.
var probed = []string{"grep", "sed", "find", "cat", "head", "tail", "wc", "sort", "uniq", "cut", "tr", "ls", "mkdir", "cp", "mv", "rm", "chmod", "xargs", "gzip", "tar", "stat", "date", "awk", "touch", "ln", "du", "basename", "dirname", "tee", "rmdir", "mktemp", "readlink", "realpath", "seq", "paste", "comm", "join", "nl", "fold", "env", "printenv", "id", "whoami", "uname", "hostname", "df", "yes", "sleep", "expand", "unexpand", "md5sum", "sha1sum", "sha224sum", "sha256sum", "sha384sum", "sha512sum", "cksum", "base64", "b2sum", "base32", "basenc", "tac", "shuf", "od", "split", "csplit", "fmt", "factor", "expr", "nproc", "timeout", "truncate", "dd", "install", "link", "unlink", "mkfifo", "nohup", "nice", "sync", "sum", "tsort", "pr", "pathchk"}

// appleSystemDirs hold macOS's own executables; coreutils such as cat live in /bin and the
// checksum tools such as sha256sum in /sbin.
var appleSystemDirs = []string{"/usr/bin", "/bin", "/sbin"}

type Selection struct {
	Implementation string   `json:"implementation"`
	Executable     string   `json:"executable,omitempty"`
	Version        string   `json:"version,omitempty"`
	Confidence     string   `json:"confidence"`
	Evidence       string   `json:"evidence"`
	Warnings       []string `json:"warnings"`
}
type Request struct {
	Name, Command, Implementation, Platform string
	NoProbe                                 bool
	// ShellBuiltin marks entries documented only as shell builtins; PATH and
	// version probes cannot describe what the user's shell runs.
	ShellBuiltin bool
	// Shells lists the shell implementations the entry documents; empty means bash.
	Shells []string
	// Single names the only implementation an entry documents (such as
	// "upstream"), so selection needs no detection.
	Single string
}

func (r Request) shells() []string {
	if len(r.Shells) == 0 {
		return []string{"bash"}
	}
	return r.Shells
}

type Detector struct{ OS, Path string }

func (d Detector) Select(ctx context.Context, r Request) (Selection, error) {
	s := Selection{Implementation: "common", Confidence: "unknown", Evidence: "Portable common reference", Warnings: []string{}}
	if r.Platform != "" && r.Implementation != "" {
		return s, fmt.Errorf("choose --platform or --implementation, not both")
	}
	if r.Implementation != "" {
		s.Implementation = r.Implementation
		s.Confidence = "explicit"
		s.Evidence = "Explicit implementation override; not a detection result"
		return s, nil
	}
	if r.Platform != "" {
		switch r.Platform {
		case "macos":
			if !r.ShellBuiltin && r.Single == "" {
				s.Implementation = "apple"
			}
		case "linux":
		case "bash":
			if !r.ShellBuiltin || !slices.Contains(r.shells(), "bash") {
				return s, fmt.Errorf("bash platform applies only to Bash builtins")
			}
			s.Implementation = "bash"
		default:
			return s, fmt.Errorf("unknown platform %q (use macos, linux, or bash)", r.Platform)
		}
		if r.Single != "" {
			s.Implementation = r.Single
		}
		s.Confidence = "target"
		s.Evidence = "Requested target platform, not local detection"
		if r.Platform == "linux" && r.Single == "" {
			s.Warnings = append(s.Warnings, "Linux does not imply GNU tools; select an implementation (for example gnu, busybox, or mawk) explicitly for its extensions.")
		}
		return s, nil
	}
	if r.Single != "" {
		s.Implementation = r.Single
		s.Confidence = "inferred"
		s.Evidence = "Entry documents a single upstream implementation; no detection needed"
		return s, nil
	}
	if r.ShellBuiltin {
		s.Evidence = "Shell builtin: parent-shell implementation cannot be inspected"
		s.Warnings = append(s.Warnings, "This entry documents "+strings.Join(r.shells(), " and ")+"; use --implementation "+strings.Join(r.shells(), " or ")+" to include shell-specific sections.")
		return s, nil
	}
	resolved, err := resolve(r.Name, d.Path)
	if err != nil {
		s.Warnings = append(s.Warnings, err.Error())
		return s, nil
	}
	s.Executable = resolved
	real, err := filepath.EvalSymlinks(resolved)
	if err != nil {
		real = resolved
	}
	if strings.Contains(filepath.Base(real), "busybox") {
		s.Implementation = "busybox"
		s.Confidence = "inferred"
		s.Evidence = "Executable resolves to a BusyBox-named binary"
	}
	if !r.NoProbe && nativeExecutable(real) && slices.Contains(probed, r.Command) {
		probeCtx, cancel := context.WithTimeout(ctx, 400*time.Millisecond)
		defer cancel()
		cmd := exec.CommandContext(probeCtx, resolved, "--version")
		cleanup := configureProbe(cmd)
		cmd.WaitDelay = 200 * time.Millisecond
		cmd.Stdin = strings.NewReader("")
		cmd.Env = append(os.Environ(), "LC_ALL=C")
		out := &limitedBuffer{}
		cmd.Stdout = out
		cmd.Stderr = out
		err := cmd.Run()
		cleanup()
		banner := out.String()
		first := strings.SplitN(banner, "\n", 2)[0]
		if probeCtx.Err() == nil && !errors.Is(err, exec.ErrWaitDelay) {
			switch {
			case d.OS == "darwin" && (strings.Contains(banner, "BSD grep") || strings.Contains(first, "-Apple") || strings.HasPrefix(first, "Apple ")):
				// Apple's banners can say "GNU compatible", "(GNU sort) ...-Apple" or "Apple gzip 479".
				s.Implementation = "apple"
			case strings.Contains(banner, "BusyBox"):
				s.Implementation = "busybox"
			case strings.Contains(first, "(GNU ") || strings.Contains(first, "(coreutils)") || strings.HasPrefix(first, "GNU Awk ") || (strings.HasPrefix(first, "gzip ") && strings.Contains(banner, "Free Software Foundation")):
				// GNU gzip prints "gzip 1.13" plus a Free Software Foundation notice, GNU dd prints
				// "dd (coreutils) 9.11" without the word GNU,
				// and gawk prints "GNU Awk 5.3.2, API 4.0, ...".
				s.Implementation = "gnu"
			case strings.HasPrefix(first, "mawk "):
				// mawk 1.3.4 prints "mawk 1.3.4 20260302".
				s.Implementation = "mawk"
			default:
				banner = ""
			}
			if banner != "" {
				s.Confidence = "detected"
				s.Evidence = "Bounded native executable --version probe"
				s.Version = firstVersion(banner)
			}
		} else {
			s.Warnings = append(s.Warnings, "Version probe or output waiting timed out; implementation may be inferred or unknown.")
		}
	} else if !r.NoProbe && !nativeExecutable(real) {
		s.Warnings = append(s.Warnings, "Version probe skipped: scripts and unknown executable formats are not executed.")
	}
	if s.Implementation == "common" && d.OS == "darwin" && isAppleSystemPath(resolved, r.Command) {
		s.Implementation = "apple"
		s.Confidence = "inferred"
		s.Evidence = "macOS system executable path (not version-certified)"
	}
	if s.Implementation == "common" {
		s.Evidence = "Executable found, but implementation is unknown"
		s.Warnings = append(s.Warnings, "Use --implementation for variant-specific help; PATH lookup cannot see parent-shell aliases or functions.")
	}
	return s, nil
}
func isAppleSystemPath(resolved, command string) bool {
	for _, dir := range appleSystemDirs {
		if resolved == dir+"/"+command {
			return true
		}
	}
	return false
}
func resolve(name, path string) (string, error) {
	candidates := []string{name}
	if !strings.ContainsRune(name, os.PathSeparator) {
		candidates = nil
		for _, dir := range filepath.SplitList(path) {
			if dir == "" {
				continue
			}
			candidates = append(candidates, filepath.Join(dir, name))
		}
	}
	for _, p := range candidates {
		info, err := os.Stat(p)
		if err == nil && !info.IsDir() && info.Mode()&0111 != 0 {
			abs, err := filepath.Abs(p)
			return abs, err
		}
	}
	return "", fmt.Errorf("%q is not an executable on PATH; showing portable help", name)
}
func nativeExecutable(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	var b [4]byte
	if _, err = io.ReadFull(f, b[:]); err != nil {
		return false
	}
	if bytes.Equal(b[:], []byte{0x7f, 'E', 'L', 'F'}) {
		return true
	}
	switch binary.BigEndian.Uint32(b[:]) {
	case 0xfeedface, 0xcefaedfe, 0xfeedfacf, 0xcffaedfe, 0xcafebabe, 0xbebafeca:
		return true
	}
	return false
}

type limitedBuffer struct{ b bytes.Buffer }

func (b *limitedBuffer) Write(p []byte) (int, error) {
	n := len(p)
	if left := 8192 - b.b.Len(); left > 0 {
		if len(p) > left {
			p = p[:left]
		}
		_, _ = b.b.Write(p)
	}
	return n, nil
}
func (b *limitedBuffer) String() string { return b.b.String() }

var version = regexp.MustCompile(`\d+\.\d+(?:\.\d+)?`)

func firstVersion(s string) string {
	for _, line := range strings.Split(s, "\n") {
		if strings.Contains(line, "BusyBox v") {
			return version.FindString(line)
		}
	}
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		s = s[:i]
	}
	return version.FindString(s)
}
