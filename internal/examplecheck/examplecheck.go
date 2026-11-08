// Package examplecheck holds the disposable fixture and expected results used to
// verify authored examples. Test helpers share it so every runner (Go tests and the
// Linux/BusyBox smoke check) judges an example the same way. It never runs commands.
package examplecheck

import (
	"fmt"
	"maps"
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

// textFiles are extra fixture files for the text-processing pages. They are not
// part of the base fixture so file-listing examples (find, ls) keep stable results.
var textFiles = map[string]string{
	"fruit.txt": "banana\napple\ncherry\napple\nbanana\napple\n",
	"nums.txt":  "10\n9\n2\n33\n1\n",
	"data.csv":  "name,qty\napple,3\nbanana,12\ncherry,7\n",
}
var textEntries = []string{"cat", "head", "tail", "wc", "sort", "uniq", "cut", "tr"}

// Setup creates the base fixture, plus any extra files the entry's examples need.
func Setup(dir, entry string) error {
	for _, sub := range Dirs {
		if err := os.Mkdir(filepath.Join(dir, sub), 0755); err != nil {
			return err
		}
	}
	files := maps.Clone(Files)
	if slices.Contains(textEntries, entry) {
		maps.Copy(files, textFiles)
	}
	for name, value := range files {
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
	// spaces compares output with runs of blanks collapsed and edges trimmed per
	// line, for tools whose column padding differs between implementations.
	spaces bool
	// notRun marks examples that block or never end (such as following a file);
	// they are reviewed but not executed.
	notRun bool
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
	// text tools
	"cat notes.txt":                             {output: "old old\nsecond\nthird\nfourth\nfifth\n"},
	"cat app.log notes.txt":                     {output: "INFO ready\nERROR broken\nerror small\nexample.org\nold old\nsecond\nthird\nfourth\nfifth\n"},
	"cat -n notes.txt":                          {output: "1 old old\n2 second\n3 third\n4 fourth\n5 fifth\n", spaces: true},
	"echo 'from stdin' | cat - notes.txt":       {output: "from stdin\nold old\nsecond\nthird\nfourth\nfifth\n"},
	"cat notes.txt > copy.txt":                  {output: ""},
	"head notes.txt":                            {output: "old old\nsecond\nthird\nfourth\nfifth\n"},
	"head -n 2 notes.txt":                       {output: "old old\nsecond\n"},
	"head -n 1 app.log logs/server.log":         {output: "==> app.log <==\nINFO ready\n\n==> logs/server.log <==\nERROR broken\n"},
	"head -c 5 notes.txt":                       {output: "old o"},
	"printf '1\\n2\\n3\\n4\\n' | head -n 2":     {output: "1\n2\n"},
	"head -n -2 notes.txt":                      {output: "old old\nsecond\nthird\n"},
	"tail notes.txt":                            {output: "old old\nsecond\nthird\nfourth\nfifth\n"},
	"tail -n 2 notes.txt":                       {output: "fourth\nfifth\n"},
	"tail -n +3 notes.txt":                      {output: "third\nfourth\nfifth\n"},
	"tail -c 6 notes.txt":                       {output: "fifth\n"},
	"tail -n 1 app.log logs/server.log":         {output: "==> app.log <==\nexample.org\n\n==> logs/server.log <==\nERROR broken\n"},
	"tail -f app.log":                           {output: "", notRun: true},
	"wc notes.txt":                              {output: "5 6 34 notes.txt", spaces: true},
	"wc -l notes.txt":                           {output: "5 notes.txt", spaces: true},
	"wc -w notes.txt":                           {output: "6 notes.txt", spaces: true},
	"wc -c notes.txt":                           {output: "34 notes.txt", spaces: true},
	"cat notes.txt | wc -l":                     {output: "5", spaces: true},
	"wc -l app.log logs/server.log":             {output: "4 app.log\n1 logs/server.log\n5 total", spaces: true},
	"wc -L notes.txt":                           {output: "7 notes.txt\n"},
	"sort fruit.txt":                            {output: "apple\napple\napple\nbanana\nbanana\ncherry\n"},
	"sort -r fruit.txt":                         {output: "cherry\nbanana\nbanana\napple\napple\napple\n"},
	"sort -n nums.txt":                          {output: "1\n2\n9\n10\n33\n"},
	"sort -u fruit.txt":                         {output: "apple\nbanana\ncherry\n"},
	"sort -t, -k2,2n data.csv":                  {output: "name,qty\napple,3\ncherry,7\nbanana,12\n"},
	"sort -rn nums.txt | head -n 1":             {output: "33\n"},
	"sort fruit.txt | uniq":                     {output: "apple\nbanana\ncherry\n"},
	"sort fruit.txt | uniq -c":                  {output: "3 apple\n2 banana\n1 cherry", spaces: true},
	"sort fruit.txt | uniq -d":                  {output: "apple\nbanana\n"},
	"sort fruit.txt | uniq -u":                  {output: "cherry\n"},
	"sort fruit.txt | uniq -c | sort -rn":       {output: "3 apple\n2 banana\n1 cherry", spaces: true},
	"uniq fruit.txt":                            {output: "banana\napple\ncherry\napple\nbanana\napple\n"},
	"cut -d, -f1 data.csv":                      {output: "name\napple\nbanana\ncherry\n"},
	"cut -d, -f2 data.csv":                      {output: "qty\n3\n12\n7\n"},
	"cut -d, -f1,2 data.csv":                    {output: "name,qty\napple,3\nbanana,12\ncherry,7\n"},
	"cut -d, -f2- data.csv":                     {output: "qty\n3\n12\n7\n"},
	"cut -c1-3 notes.txt":                       {output: "old\nsec\nthi\nfou\nfif\n"},
	"echo 'a b c' | cut -d' ' -f2":              {output: "b\n"},
	"cut --complement -d, -f1 data.csv":         {output: "qty\n3\n12\n7\n"},
	"echo 'hello' | tr 'a-z' 'A-Z'":             {output: "HELLO\n"},
	"echo 'hello' | tr '[:lower:]' '[:upper:]'": {output: "HELLO\n"},
	"echo 'a   b    c' | tr -s ' '":             {output: "a b c\n"},
	"echo 'hello world' | tr -d 'o'":            {output: "hell wrld\n"},
	"echo 'abc123' | tr -cd '0-9'":              {output: "123"},
	"echo 'a:b:c' | tr ':' '\\n'":               {output: "a\nb\nc\n"},
	"tr ',' '\\t' < data.csv":                   {output: "name\tqty\napple\t3\nbanana\t12\ncherry\t7\n"},
	// File tools: each example ends by listing or reading the result.
	"ls":                             {output: "app.log\ncache\nlogs\nnotes.txt\nprojects\n"},
	"ls logs":                        {output: "server.log\n"},
	"ls *.log":                       {output: "app.log\n"},
	"ls -d */":                       {output: "cache/\nlogs/\nprojects/\n"},
	"ls -a logs":                     {output: ".\n..\nserver.log\n"},
	"ls -R logs projects":            {output: "logs:\nserver.log\n\nprojects:\n"},
	"ls -r":                          {output: "projects\nnotes.txt\nlogs\ncache\napp.log\n"},
	"ls -l notes.txt | cut -c1-10":   {output: "-rw-r--r--\n"},
	"ls --group-directories-first":   {output: "cache\nlogs\nprojects\napp.log\nnotes.txt\n"},
	"mkdir reports && ls -d reports": {output: "reports\n"},
	"mkdir -p reports/2026/q1 && ls reports/2026":         {output: "q1\n"},
	"mkdir one two && ls -d one two":                      {output: "one\ntwo\n"},
	"mkdir -m 700 private && ls -ld private | cut -c1-10": {output: "drwx------\n"},
	"mkdir -p cache && echo ok":                           {output: "ok\n"},
	"mkdir -pv reports/2026/q2":                           {output: "mkdir: created directory 'reports'\nmkdir: created directory 'reports/2026'\nmkdir: created directory 'reports/2026/q2'\n"},
	"cp notes.txt copy.txt && cat copy.txt":               {output: "old old\nsecond\nthird\nfourth\nfifth\n"},
	"cp app.log logs/ && ls logs":                         {output: "app.log\nserver.log\n"},
	"cp notes.txt app.log logs && ls logs":                {output: "app.log\nnotes.txt\nserver.log\n"},
	"cp -R logs logs-copy && ls logs-copy":                {output: "server.log\n"},
	"cp -p notes.txt kept.txt && ls kept.txt":             {output: "kept.txt\n"},
	"cp -t logs notes.txt app.log && ls logs":             {output: "app.log\nnotes.txt\nserver.log\n"},
	"mv notes.txt renamed.txt && head -n 1 renamed.txt":   {output: "old old\n"},
	"mv app.log logs/ && ls logs":                         {output: "app.log\nserver.log\n"},
	"mv app.log notes.txt logs && ls logs":                {output: "app.log\nnotes.txt\nserver.log\n"},
	"mv logs old-logs && ls old-logs":                     {output: "server.log\n"},
	"mv -f notes.txt app.log && head -n 1 app.log":        {output: "old old\n"},
	"mv -t logs notes.txt app.log && ls logs":             {output: "app.log\nnotes.txt\nserver.log\n"},
	"rm app.log && ls":                                    {output: "cache\nlogs\nnotes.txt\nprojects\n"},
	"rm app.log notes.txt && ls":                          {output: "cache\nlogs\nprojects\n"},
	"rm -r logs && ls":                                    {output: "app.log\ncache\nnotes.txt\nprojects\n"},
	"rm -f missing.txt && echo done":                      {output: "done\n"},
	"rm -I app.log && ls":                                 {output: "cache\nlogs\nnotes.txt\nprojects\n"},
	// cd: only success is checked; the printed directory is a temporary path.
	"cd projects":        {},
	"cd ..":              {},
	"cd":                 {},
	"cd ~/projects":      {},
	"cd projects && pwd": {},
	"cd -":               {},
	"cd -P projects":     {},
}

// NotRun reports whether code is reviewed only and must not be executed.
func NotRun(code string) bool { return expected[strings.TrimSpace(code)].notRun }

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
	if want.notRun || strings.HasPrefix(code, "cd") {
		return nil
	}
	if want.spaces {
		if collapse(output) != collapse(want.output) {
			return fmt.Errorf("output %q, want %q after collapsing blanks", output, want.output)
		}
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

func collapse(s string) string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	for i, line := range lines {
		lines[i] = strings.Join(strings.Fields(line), " ")
	}
	return strings.Join(lines, "\n")
}

func sortedLines(s string) []string {
	lines := strings.Split(strings.TrimSuffix(s, "\n"), "\n")
	slices.Sort(lines)
	return lines
}
