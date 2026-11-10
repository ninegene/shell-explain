package examplecheck

import (
	"fmt"
	"regexp"
	"strings"
)

// TryItBlock is one command block from a guide's "Try it" section. When a
// fenced text block follows it directly, that block is the exact output the
// command must print; otherwise only a successful exit is checked, because
// the output (disk sizes, process numbers) differs on every computer.
type TryItBlock struct {
	Code      string
	Expect    string
	HasExpect bool
}

var fence = regexp.MustCompile("(?ms)^```([a-z]*)\n(.*?)^```$")

// TryItBlocks extracts the sh blocks of a Markdown section, each paired with
// the text block that follows it, if any. It never runs anything.
func TryItBlocks(markdown string) []TryItBlock {
	var blocks []TryItBlock
	matches := fence.FindAllStringSubmatchIndex(markdown, -1)
	for i, m := range matches {
		lang, body := markdown[m[2]:m[3]], markdown[m[4]:m[5]]
		if lang != "sh" {
			continue
		}
		block := TryItBlock{Code: body}
		if i+1 < len(matches) {
			next := matches[i+1]
			between := markdown[m[1]:next[0]]
			// The output block follows its command with no other block between.
			if markdown[next[2]:next[3]] == "text" && !strings.Contains(between, "```") {
				block.Expect, block.HasExpect = markdown[next[4]:next[5]], true
			}
		}
		blocks = append(blocks, block)
	}
	return blocks
}

// unsafeTryIt names commands a beginner guide must never ask a reader to run.
var unsafeTryIt = regexp.MustCompile(`(^|[^a-zA-Z0-9_-])(sudo|su|rm|rmdir|kill|killall|pkill|shutdown|reboot|halt|poweroff|mkfs|fdisk|dd|chown|chgrp|ln|wget|ssh|scp)([^a-zA-Z0-9_-]|$)|>\s*/(dev|etc|usr|bin|sbin|var|System|Library)\b`)

// UnsafeTryIt returns the first thing in code that a "Try it" block must not
// contain, or "" when it is safe: only reading, and writing inside the current
// folder.
func UnsafeTryIt(code string) string {
	// A backslash at the end of a line continues the command, so check the
	// joined command: the words after it belong to the same mv.
	code = strings.ReplaceAll(code, "\\\n", " ")
	hasMove := false
	for _, line := range strings.Split(code, "\n") {
		if curlWord.MatchString(line) && !curlStaysLocal(line) {
			return "curl"
		}
		if m := unsafeTryIt.FindStringSubmatch(line); m != nil {
			if m[2] != "" {
				return m[2]
			}
			return "a redirect to a system path"
		}
		if moveWord.MatchString(line) {
			hasMove = true
			if movesOutside(line) {
				return "mv outside the current folder"
			}
		}
	}
	if hasMove {
		// A cd must go to a plain folder name inside the current folder: a bare
		// cd goes home and "cd -" goes back, which would make a relative mv
		// unsafe.
		for _, m := range cdTarget.FindAllStringSubmatch(code, -1) {
			if !plainFolder.MatchString(strings.TrimSpace(m[1])) {
				return "mv outside the current folder"
			}
		}
	}
	return ""
}

// curl is allowed in a Try it block only to read a file:// URL, so the tests
// never touch the network: the line must name a file:// URL, no other URL
// scheme, and no path outside the current folder.
var (
	curlWord  = regexp.MustCompile(`(^|[^a-zA-Z0-9_-])curl([^a-zA-Z0-9_-]|$)`)
	urlScheme = regexp.MustCompile(`(?i)[a-z][a-z0-9+.-]*://`)
)

// UsesCurl reports whether code runs curl, so a machine without curl can skip
// the block instead of failing it.
func UsesCurl(code string) bool { return curlWord.MatchString(code) }

func curlStaysLocal(line string) bool {
	if !strings.Contains(line, "file://") || strings.ContainsAny(line, "`") || strings.Contains(line, "$(") {
		return false
	}
	for _, scheme := range urlScheme.FindAllString(line, -1) {
		if !strings.EqualFold(scheme, "file://") {
			return false
		}
	}
	// "$PWD" is the one variable allowed: it keeps the URL in this folder.
	rest := strings.ReplaceAll(line, "file://$PWD/", "file://")
	if strings.Contains(rest, "file:///") || strings.Contains(rest, "$") {
		return false
	}
	if leavesFolder.MatchString(strings.ReplaceAll(rest, "file://", "file:")) {
		return false
	}
	return curlTargetsAreFiles(line)
}

// curlTakesValue lists the options whose next word is a value, not a URL.
var curlTakesValue = map[string]bool{"-o": true, "-w": true, "-H": true, "-X": true, "-d": true, "-A": true, "-e": true, "--output": true, "--write-out": true, "--header": true, "--request": true, "--data": true}

// curlTargetsAreFiles rejects a bare host name such as "example.org", which
// curl would fetch over the network: every word of the curl command that is
// not an option or an option's value must be a file:// URL.
func curlTargetsAreFiles(line string) bool {
	line = regexp.MustCompile(`&&|\|\||;`).ReplaceAllString(line, "\n")
	for _, segment := range strings.Split(line, "\n") {
		words := strings.Fields(segment)
		if len(words) == 0 || strings.Trim(words[0], "\"'") != "curl" {
			continue
		}
		skip := false
		for _, w := range words[1:] {
			switch {
			case skip:
				skip = false
			case strings.HasPrefix(w, "-"):
				skip = curlTakesValue[w]
			case !strings.Contains(w, "file://"):
				return false
			}
		}
	}
	return true
}

// mv is allowed because the guide on renaming and moving files needs it, but
// only on files the block itself made in the current folder. A line with mv is
// rejected when any word starts with a slash or a tilde, when ".." appears
// anywhere, or when it uses a backtick, a command substitution or a variable.
// The one variable allowed is the loop variable of a for line, so the
// rename-many example works.
var (
	moveWord = regexp.MustCompile(`(^|[^a-zA-Z0-9_-])mv([^a-zA-Z0-9_-]|$)`)
	cdTarget = regexp.MustCompile(`(?m)(?:^|[^a-zA-Z0-9_-])cd(?:[ 	]+([^;&|
]*))?[ 	]*(?:[;&|
]|$)`)
	plainFolder  = regexp.MustCompile(`^[A-Za-z0-9_][A-Za-z0-9_.-]*(/[A-Za-z0-9_][A-Za-z0-9_.-]*)*/?$`)
	leavesFolder = regexp.MustCompile(`(^|\s)["']?(/|~)|\.\.`)
	loopLine     = regexp.MustCompile(`^\s*for\s+([A-Za-z_][A-Za-z0-9_]*)\s+in\b`)
	variableRef  = regexp.MustCompile(`\$\{?([A-Za-z_][A-Za-z0-9_]*)|\$[^A-Za-z_{]|\$\{[^A-Za-z_]`)
)

func movesOutside(line string) bool {
	if leavesFolder.MatchString(line) || strings.ContainsAny(line, "`") || strings.Contains(line, "$(") {
		return true
	}
	loop := ""
	if m := loopLine.FindStringSubmatch(line); m != nil {
		loop = m[1]
	}
	for _, m := range variableRef.FindAllStringSubmatch(line, -1) {
		if m[1] == "" || m[1] != loop {
			return true
		}
	}
	return false
}

// CheckTryIt compares a command's output with the output the guide shows.
func CheckTryIt(b TryItBlock, output string) error {
	if b.HasExpect && output != b.Expect {
		return fmt.Errorf("output %q does not match the %q shown in the guide", output, b.Expect)
	}
	return nil
}
