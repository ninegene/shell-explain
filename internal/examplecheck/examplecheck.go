// Package examplecheck holds the disposable fixture and expected results used to
// verify authored examples. Test helpers share it so every runner (Go tests and the
// Linux/BusyBox smoke check) judges an example the same way. It never runs commands.
package examplecheck

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Dirs and Files describe the synthetic working directory examples run in.
var Dirs = []string{"projects", "logs", "cache"}
var Files = map[string]string{
	"app.log":         "INFO ready\nERROR broken\nerror small\nexample.org\n",
	"logs/server.log": "ERROR broken\n",
	"notes.txt":       "old old\nsecond\nthird\nfourth\nfifth\n",
}

// Setup creates the fixture inside dir.
func Setup(dir string) error {
	for _, sub := range Dirs {
		if err := os.Mkdir(filepath.Join(dir, sub), 0755); err != nil {
			return err
		}
	}
	for name, value := range Files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(value), 0644); err != nil {
			return err
		}
	}
	return nil
}

// expectation is the exact output of an example. Unordered output (directory
// traversal order) is compared as a sorted set of lines.
type expectation struct {
	output    string
	unordered bool
}

var expected = map[string]expectation{
	// grep
	"grep 'ERROR' app.log":                     {output: "ERROR broken\n"},
	"grep -n -i 'error' app.log":               {output: "2:ERROR broken\n3:error small\n"},
	"grep -F 'example.org' app.log":            {output: "example.org\n"},
	"grep -v 'ERROR' app.log":                  {output: "INFO ready\nerror small\nexample.org\n"},
	"grep -c 'ERROR' app.log":                  {output: "1\n"},
	"grep -l 'ERROR' app.log logs/server.log":  {output: "app.log\nlogs/server.log\n"},
	"grep -w 'error' app.log":                  {output: "error small\n"},
	"grep -E 'ERROR|INFO' app.log":             {output: "INFO ready\nERROR broken\n"},
	"grep -e 'INFO' -e 'example' app.log":      {output: "INFO ready\nexample.org\n"},
	"grep -C 1 'ERROR' app.log":                {output: "INFO ready\nERROR broken\nerror small\n"},
	"grep -r 'ERROR' logs":                     {output: "logs/server.log:ERROR broken\n"},
	"grep -rn --include='*.log' 'ERROR' logs/": {output: "logs/server.log:1:ERROR broken\n"},
	// sed (in-place examples are checked against files by Verify)
	"sed 's/old/new/' notes.txt":                       {output: "new old\nsecond\nthird\nfourth\nfifth\n"},
	"sed 's/old/new/g' notes.txt":                      {output: "new new\nsecond\nthird\nfourth\nfifth\n"},
	"sed -n '2,4p' notes.txt":                          {output: "second\nthird\nfourth\n"},
	"sed '/second/d' notes.txt":                        {output: "old old\nthird\nfourth\nfifth\n"},
	"sed -n '/third/p' notes.txt":                      {output: "third\n"},
	"sed '1d' notes.txt":                               {output: "second\nthird\nfourth\nfifth\n"},
	"sed -n '$p' notes.txt":                            {output: "fifth\n"},
	"sed -e 's/old/new/' -e 's/second/2nd/' notes.txt": {output: "new old\n2nd\nthird\nfourth\nfifth\n"},
	"echo '/usr/local/bin' | sed 's#/usr/local#/opt#'": {output: "/opt/bin\n"},
	"echo 'a1b22' | sed -E 's/[0-9]+/N/g'":             {output: "aNbN\n"},
	"sed -i.bak 's/old/new/g' notes.txt":               {},
	"sed -i '.bak' 's/old/new/g' notes.txt":            {},
	// find
	"find . -type f -name '*.log'":                            {output: "./app.log\n./logs/server.log\n", unordered: true},
	"find . -type d -name 'cache'":                            {output: "./cache\n"},
	"find . -type f -iname '*.LOG'":                           {output: "./app.log\n./logs/server.log\n", unordered: true},
	"find . -type f ! -name '*.log'":                          {output: "./notes.txt\n"},
	"find . -type f \\( -name '*.log' -o -name '*.txt' \\)":   {output: "./app.log\n./logs/server.log\n./notes.txt\n", unordered: true},
	"find . -type f -mtime -7":                                {output: "./app.log\n./logs/server.log\n./notes.txt\n", unordered: true},
	"find . -type d":                                          {output: ".\n./cache\n./logs\n./projects\n", unordered: true},
	"find . -type f -name '*.log' -exec grep -l 'ERROR' {} +": {output: "./app.log\n./logs/server.log\n", unordered: true},
	"find . -maxdepth 2 -type f -name '*.log'":                {output: "./app.log\n./logs/server.log\n", unordered: true},
	// cd: only success is checked; the printed directory is a temporary path.
	"cd projects":        {},
	"cd ..":              {},
	"cd":                 {},
	"cd ~/projects":      {},
	"cd projects && pwd": {},
	"cd -":               {},
	"cd -P projects":     {},
}

// HasExpectation reports whether code has a recorded expectation.
func HasExpectation(code string) bool {
	_, ok := expected[strings.TrimSpace(code)]
	return ok
}

// Verify checks combined output and the fixture directory after an example ran
// in dir. It fails for an example with no recorded expectation, so new examples
// cannot ship unchecked.
func Verify(code, dir, output string) error {
	code = strings.TrimSpace(code)
	want, ok := expected[code]
	if !ok {
		return fmt.Errorf("no expected result recorded for example %q", code)
	}
	if strings.Contains(code, "sed -i") {
		got, err := os.ReadFile(filepath.Join(dir, "notes.txt"))
		if err != nil || !strings.HasPrefix(string(got), "new new") {
			return fmt.Errorf("in-place edit result %q: %v", got, err)
		}
		backup, err := os.ReadFile(filepath.Join(dir, "notes.txt.bak"))
		if err != nil || !strings.HasPrefix(string(backup), "old old") {
			return fmt.Errorf("backup not preserved: %v", err)
		}
		return nil
	}
	if strings.HasPrefix(code, "cd") {
		return nil
	}
	if want.unordered {
		if !slices.Equal(sortedLines(output), sortedLines(want.output)) {
			return fmt.Errorf("output %q, want the lines of %q in any order", output, want.output)
		}
		return nil
	}
	if output != want.output {
		return fmt.Errorf("output %q, want %q", output, want.output)
	}
	return nil
}

func sortedLines(s string) []string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	slices.Sort(lines)
	return lines
}
