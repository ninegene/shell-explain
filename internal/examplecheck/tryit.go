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
var unsafeTryIt = regexp.MustCompile(`(^|[^a-zA-Z0-9_-])(sudo|su|rm|rmdir|kill|killall|pkill|shutdown|reboot|halt|poweroff|mkfs|fdisk|dd|chown|chgrp|ln|curl|wget|ssh|scp)([^a-zA-Z0-9_-]|$)|>\s*/(dev|etc|usr|bin|sbin|var|System|Library)\b`)

// UnsafeTryIt returns the first thing in code that a "Try it" block must not
// contain, or "" when it is safe: only reading, and writing inside the current
// folder.
func UnsafeTryIt(code string) string {
	for _, line := range strings.Split(code, "\n") {
		if m := unsafeTryIt.FindStringSubmatch(line); m != nil {
			if m[2] != "" {
				return m[2]
			}
			return "a redirect to a system path"
		}
		if moveWord.MatchString(line) && leavesFolder.MatchString(line) {
			return "mv outside the current folder"
		}
	}
	return ""
}

// mv is allowed because the guide on renaming and moving files needs it, but
// only on relative paths: a word that starts with a slash, a tilde or two dots
// could reach the reader's own files.
var (
	moveWord     = regexp.MustCompile(`(^|[^a-zA-Z0-9_-])mv([^a-zA-Z0-9_-]|$)`)
	leavesFolder = regexp.MustCompile(`(^|\s)["']?(/|~|\.\.)`)
)

// CheckTryIt compares a command's output with the output the guide shows.
func CheckTryIt(b TryItBlock, output string) error {
	if b.HasExpect && output != b.Expect {
		return fmt.Errorf("output %q does not match the %q shown in the guide", output, b.Expect)
	}
	return nil
}
