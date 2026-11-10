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

// Companions lists programs that a command runs by name, which must come from the same
// implementation as the command: GNU diff3 and sdiff call diff.
var Companions = map[string][]string{"diff3": {"diff"}, "sdiff": {"diff"}}

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
var textEntries = []string{"cat", "head", "tail", "wc", "sort", "uniq", "cut", "tr", "awk"}

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
	// chmod, xargs, gzip, tar: examples end by showing the result.
	"chmod 600 notes.txt && ls -l notes.txt | cut -c1-10":                                      {output: "-rw-------\n"},
	"chmod u+x notes.txt && ls -l notes.txt | cut -c1-10":                                      {output: "-rwxr--r--\n"},
	"chmod a+x notes.txt && ls -l notes.txt | cut -c1-10":                                      {output: "-rwxr-xr-x\n"},
	"chmod go-r app.log && ls -l app.log | cut -c1-10":                                         {output: "-rw-------\n"},
	"chmod u=rw,go=r notes.txt && ls -l notes.txt | cut -c1-10":                                {output: "-rw-r--r--\n"},
	"chmod -R 700 logs && ls -ld logs logs/server.log | cut -c1-10":                            {output: "drwx------\n-rwx------\n"},
	"chmod --reference=app.log notes.txt && ls -l notes.txt | cut -c1-10":                      {output: "-rw-r--r--\n"},
	"printf 'a\\nb\\n' | xargs echo":                                                           {output: "a b\n"},
	"printf 'a\\nb\\nc\\n' | xargs -n 2 echo":                                                  {output: "a b\nc\n"},
	"printf 'one\\ntwo\\n' | xargs -I {} echo item-{}":                                         {output: "item-one\nitem-two\n"},
	"printf 'a b\\0c\\0' | xargs -0 -n 1 echo":                                                 {output: "a b\nc\n"},
	"printf 'app.log\\nnotes.txt\\n' | xargs wc -l":                                            {output: "4 app.log\n5 notes.txt\n9 total", spaces: true},
	"printf '' | xargs -r echo ran; echo finished":                                             {output: "finished\n"},
	"printf 'a,b,c' | xargs -d , -n 1 echo":                                                    {output: "a\nb\nc\n"},
	"gzip app.log && ls":                                                                       {output: "app.log.gz\ncache\nlogs\nnotes.txt\nprojects\n"},
	"gzip -k app.log && ls app.log*":                                                           {output: "app.log\napp.log.gz\n"},
	"gzip -k app.log && gzip -dc app.log.gz":                                                   {output: "INFO ready\nERROR broken\nerror small\nexample.org\n"},
	"gzip notes.txt && gzip -d notes.txt.gz && head -n 1 notes.txt":                            {output: "old old\n"},
	"gzip -c app.log > copy.log.gz && ls *.gz":                                                 {output: "copy.log.gz\n"},
	"gzip -k -S .z app.log && ls app.log*":                                                     {output: "app.log\napp.log.z\n"},
	"tar -cf logs.tar logs && tar -tf logs.tar":                                                {output: "logs/\nlogs/server.log\n"},
	"tar -cf backup.tar app.log notes.txt && tar -tf backup.tar":                               {output: "app.log\nnotes.txt\n"},
	"tar -czf logs.tgz logs && tar -tzf logs.tgz":                                              {output: "logs/\nlogs/server.log\n"},
	"mkdir restore && tar -cf logs.tar logs && tar -xf logs.tar -C restore && ls restore/logs": {output: "server.log\n"},
	"tar -cf logs.tar logs && tar -xOf logs.tar logs/server.log":                               {output: "ERROR broken\n"},
	"tar --exclude='*.log' -cf logs.tar logs && tar -tf logs.tar":                              {output: "logs/\n"},
	"tar -cf logs.tar --transform 's,^logs,archive,' logs && tar -tf logs.tar":                 {output: "archive/\narchive/server.log\n"},
	// date: the clock changes, so current-time examples check only the shape.
	"date +%Y-%m-%d | grep -c '^[0-9]\\{4\\}-[0-9]\\{2\\}-[0-9]\\{2\\}$'": {output: "1\n"},
	"date +%H:%M:%S | grep -c '^[0-9][0-9]:[0-9][0-9]:[0-9][0-9]$'":       {output: "1\n"},
	"date +%s | grep -c '^[0-9]*$'":                                       {output: "1\n"},
	"date -u +%Y-%m-%dT%H:%M:%SZ | grep -c 'Z$'":                          {output: "1\n"},
	"date -u +%Z":                                                 {output: "UTC\n"},
	"date -u -d @86400 +%Y-%m-%d":                                 {output: "1970-01-02\n"},
	"date -u -d '2026-11-08 12:00' +%s":                           {output: "1794139200\n"},
	"date -u -r 86400 +%Y-%m-%d":                                  {output: "1970-01-02\n"},
	"date -u -j -f '%Y-%m-%d %H:%M:%S' '2026-11-08 12:00:00' +%s": {output: "1794139200\n"},
	// stat: the report formats differ, so only exit status and chosen fields are checked.
	"stat notes.txt > /dev/null && echo exists":     {output: "exists\n"},
	"stat missing.txt 2> /dev/null || echo missing": {output: "missing\n"},
	"stat -c '%s' notes.txt":                        {output: "34\n"},
	"stat -c '%a %n' notes.txt":                     {output: "644 notes.txt\n"},
	"stat -c '%F' logs":                             {output: "directory\n"},
	"stat -f '%z' notes.txt":                        {output: "34\n"},
	"stat -f '%Lp %N' notes.txt":                    {output: "644 notes.txt\n"},
	"stat -f '%HT' logs":                            {output: "Directory\n"},
	// awk: the shared examples run on macOS awk, gawk, mawk and BusyBox awk. Array
	// order is unspecified, so that example is compared as a set of lines.
	"awk '{print $1}' notes.txt":                                             {output: "old\nsecond\nthird\nfourth\nfifth\n"},
	"awk -F, '{print $1}' data.csv":                                          {output: "name\napple\nbanana\ncherry\n"},
	"awk -F, -v min=5 'NR > 1 && $2 > min {print $1}' data.csv":              {output: "banana\ncherry\n"},
	"awk -F, 'NR > 1 {sum += $2} END {print sum}' data.csv":                  {output: "22\n"},
	"awk '/ERROR/ {print $2}' app.log":                                       {output: "broken\n"},
	"awk '{print NF}' notes.txt":                                             {output: "2\n1\n1\n1\n1\n"},
	"awk -F, 'NR > 1 {printf \"%-8s %3d\\n\", $1, $2}' data.csv":             {output: "apple      3\nbanana    12\ncherry     7\n"},
	"awk '{count[$1]++} END {for (w in count) print w, count[w]}' fruit.txt": {output: "apple 3\nbanana 2\ncherry 1\n", unordered: true},
	"awk 'BEGIN {PROCINFO[\"sorted_in\"] = \"@ind_str_asc\"} {count[$1]++} END {for (w in count) print w, count[w]}' fruit.txt": {output: "apple 3\nbanana 2\ncherry 1\n"},
	"echo 'a1b22' | awk '{print gensub(/[0-9]+/, \"N\", \"g\")}'":                                                               {output: "aNbN\n"},
	"awk 'BEGIN {print strftime(\"%Y-%m-%d\", 86400, 1)}'":                                                                      {output: "1970-01-02\n"},
	// touch, ln, du, basename, dirname, tee: results use the base fixture. du prints
	// disk usage that depends on the file system, so only the shape of a line is checked.
	"touch new.txt && ls new.txt":                                                                                     {output: "new.txt\n"},
	"touch a.txt b.txt && ls a.txt b.txt":                                                                             {output: "a.txt\nb.txt\n"},
	"touch notes.txt && head -n 1 notes.txt":                                                                          {output: "old old\n"},
	"touch -c missing.txt && ls missing.txt 2> /dev/null || echo not created":                                         {output: "not created\n"},
	"touch -t 200101020304 old.txt && test notes.txt -nt old.txt && echo notes is newer":                              {output: "notes is newer\n"},
	"touch -r notes.txt copy.txt && test ! copy.txt -nt notes.txt && test ! notes.txt -nt copy.txt && echo same time": {output: "same time\n"},
	"touch -d '2001-01-02 03:04' old.txt && test notes.txt -nt old.txt && echo notes is newer":                        {output: "notes is newer\n"},
	"ln -s notes.txt link.txt && head -n 1 link.txt":                                                                  {output: "old old\n"},
	"ln -s notes.txt link.txt && test -L link.txt && echo symlink":                                                    {output: "symlink\n"},
	"ln notes.txt copy.txt && echo extra >> copy.txt && tail -n 1 notes.txt":                                          {output: "extra\n"},
	"ln -s ../notes.txt projects/notes.txt && head -n 1 projects/notes.txt":                                           {output: "old old\n"},
	"ln -s notes.txt current && ln -sf app.log current && head -n 1 current":                                          {output: "INFO ready\n"},
	"ln -s logs latest && ls latest":                                                                                  {output: "server.log\n"},
	"ln -s missing.txt broken && { test -e broken || echo broken link; }":                                             {output: "broken link\n"},
	"ln -sr notes.txt projects/notes.txt && head -n 1 projects/notes.txt":                                             {output: "old old\n"},
	"du -sk logs | cut -f2": {output: "logs\n"},
	"du -sh logs | grep -c '^[0-9.]*[BKMGT]*[[:space:]]logs$'": {output: "1\n"},
	"du -a logs | cut -f2":                                               {output: "logs\nlogs/server.log\n", unordered: true},
	"du -d 1 . | cut -f2":                                                {output: ".\n./cache\n./logs\n./projects\n", unordered: true},
	"du -c notes.txt app.log | tail -n 1 | cut -f2":                      {output: "total\n"},
	"du -sk logs | grep -c '^[0-9][0-9]*[[:space:]]logs$'":               {output: "1\n"},
	"du -b notes.txt":                                                    {output: "34\tnotes.txt\n"},
	"basename logs/server.log":                                           {output: "server.log\n"},
	"basename logs/server.log .log":                                      {output: "server\n"},
	"basename /usr/local/bin/":                                           {output: "bin\n"},
	"for f in logs/*.log; do echo \"$(basename \"$f\" .log).txt\"; done": {output: "server.txt\n"},
	"basename -a logs/server.log projects/notes.txt":                     {output: "server.log\nnotes.txt\n"},
	"dirname logs/server.log":                                            {output: "logs\n"},
	"dirname server.log":                                                 {output: ".\n"},
	"dirname /usr/local/bin/tool":                                        {output: "/usr/local/bin\n"},
	"dirname /usr/local/bin/":                                            {output: "/usr/local\n"},
	"mkdir -p \"$(dirname cache/data/out.txt)\" && ls cache":             {output: "data\n"},
	"dirname logs/server.log projects/notes.txt":                         {output: "logs\nprojects\n"},
	"echo hello | tee out.txt && cat out.txt":                            {output: "hello\nhello\n"},
	"echo hello | tee out.txt > /dev/null && cat out.txt":                {output: "hello\n"},
	"echo one | tee log.txt > /dev/null && echo two | tee -a log.txt > /dev/null && cat log.txt": {output: "one\ntwo\n"},
	"echo hi | tee a.txt b.txt > /dev/null && cat a.txt b.txt":                                   {output: "hi\nhi\n"},
	"printf 'b\\na\\n' | tee unsorted.txt | sort && cat unsorted.txt":                            {output: "a\nb\nb\na\n"},
	"echo hello | tee -p out.txt > /dev/null && cat out.txt":                                     {output: "hello\n"},
	// rmdir, mktemp, readlink, realpath, seq, paste, comm, join, nl, fold. mktemp
	// and realpath print temporary paths, so examples keep only a stable part.
	"mkdir empty && rmdir empty && echo removed":                                 {output: "removed\n"},
	"mkdir one two && rmdir one two && echo ok":                                  {output: "ok\n"},
	"rmdir logs 2> /dev/null || echo not empty":                                  {output: "not empty\n"},
	"mkdir -p a/b/c && rmdir -p a/b/c && { ls -d a 2> /dev/null || echo gone; }": {output: "gone\n"},
	"rmdir --ignore-fail-on-non-empty logs && ls logs":                           {output: "server.log\n"},
	"f=$(mktemp) && test -f \"$f\" && echo created; rm -f \"$f\"":                {output: "created\n"},
	"d=$(mktemp -d) && test -d \"$d\" && echo created; rmdir \"$d\"":             {output: "created\n"},
	"f=$(mktemp cache/item.XXXXXX) && ls cache | grep -c .":                      {output: "1\n"},
	"f=$(mktemp work.XXXXXX) && echo \"${f%.*}\"":                                {output: "work\n"},
	"f=$(mktemp item.XXXXXX --suffix=.txt) && echo \"${f##*.}\"":                 {output: "txt\n"},
	"ln -s notes.txt link && readlink link":                                      {output: "notes.txt\n"},
	"ln -s notes.txt link && basename \"$(readlink -f link)\"":                   {output: "notes.txt\n"},
	"readlink notes.txt || echo not a link":                                      {output: "not a link\n"},
	"ln -s ../notes.txt projects/up && readlink projects/up":                     {output: "../notes.txt\n"},
	"readlink -e missing.txt || echo missing":                                    {output: "missing\n"},
	"realpath logs/../notes.txt | grep -c '^/'":                                  {output: "1\n"},
	"basename \"$(realpath logs/../notes.txt)\"":                                 {output: "notes.txt\n"},
	"ln -s notes.txt link && basename \"$(realpath link)\"":                      {output: "notes.txt\n"},
	"cd logs && basename \"$(realpath .)\"":                                      {output: "logs\n"},
	"realpath --relative-to=logs notes.txt":                                      {output: "../notes.txt\n"},
	"basename \"$(realpath -m logs/new.txt)\"":                                   {output: "new.txt\n"},
	"seq 3":                     {output: "1\n2\n3\n"},
	"seq 2 4":                   {output: "2\n3\n4\n"},
	"seq 1 2 7":                 {output: "1\n3\n5\n7\n"},
	"seq 5 -1 1 | paste -sd, -": {output: "5,4,3,2,1\n"},
	"seq -w 8 10":               {output: "08\n09\n10\n"},
	"for i in $(seq 3); do echo \"run $i\"; done": {output: "run 1\nrun 2\nrun 3\n"},
	"seq -s, 1 4": {output: "1,2,3,4\n"},
	"printf 'a\\nb\\n' > names.txt && printf '1\\n2\\n' > nums2.txt && paste names.txt nums2.txt":     {output: "a\t1\nb\t2\n"},
	"printf 'a\\nb\\n' > names.txt && printf '1\\n2\\n' > nums2.txt && paste -d, names.txt nums2.txt": {output: "a,1\nb,2\n"},
	"printf 'a\\nb\\nc\\n' > names.txt && paste -s names.txt":                                         {output: "a\tb\tc\n"},
	"printf 'a\\nb\\nc\\n' > names.txt && paste -sd, names.txt":                                       {output: "a,b,c\n"},
	"seq 6 | paste - - -":  {output: "1\t2\t3\n4\t5\t6\n"},
	"seq 3 | paste -sd+ -": {output: "1+2+3\n"},
	"printf 'a\\nb\\nc\\n' > one.txt && printf 'b\\nc\\nd\\n' > two.txt && comm -12 one.txt two.txt":                                            {output: "b\nc\n"},
	"printf 'a\\nb\\nc\\n' > one.txt && printf 'b\\nc\\nd\\n' > two.txt && comm -23 one.txt two.txt":                                            {output: "a\n"},
	"printf 'a\\nb\\nc\\n' > one.txt && printf 'b\\nc\\nd\\n' > two.txt && comm -13 one.txt two.txt":                                            {output: "d\n"},
	"printf 'c\\na\\nb\\n' | sort > one.txt && printf 'd\\nb\\nc\\n' | sort > two.txt && comm -12 one.txt two.txt":                              {output: "b\nc\n"},
	"printf 'b\\na\\n' > one.txt && printf 'a\\nb\\n' > two.txt && { comm --check-order one.txt two.txt > /dev/null 2>&1 || echo not sorted; }": {output: "not sorted\n"},
	"printf '1 alice\\n2 bob\\n3 carol\\n' > names.txt && printf '1 admin\\n3 user\\n' > roles.txt && join names.txt roles.txt":                 {output: "1 alice admin\n3 carol user\n"},
	"printf '1 alice\\n2 bob\\n3 carol\\n' > names.txt && printf '1 admin\\n3 user\\n' > roles.txt && join -a 1 names.txt roles.txt":            {output: "1 alice admin\n2 bob\n3 carol user\n"},
	"printf '1 alice\\n2 bob\\n3 carol\\n' > names.txt && printf '1 admin\\n3 user\\n' > roles.txt && join -v 1 names.txt roles.txt":            {output: "2 bob\n"},
	"printf '1,alice\\n2,bob\\n' > names.csv && printf '1,admin\\n2,user\\n' > roles.csv && join -t, names.csv roles.csv":                       {output: "1,alice,admin\n2,bob,user\n"},
	"printf '1 alice\\n2 bob\\n' > names.txt && printf 'admin 1\\nuser 2\\n' > roles2.txt && join -2 2 names.txt roles2.txt":                    {output: "1 alice admin\n2 bob user\n"},
	"printf 'id name\\n1 alice\\n' > names.txt && printf 'id role\\n1 admin\\n' > roles.txt && join --header names.txt roles.txt":               {output: "id name role\n1 alice admin\n"},
	"nl notes.txt":                                           {output: "1 old old\n2 second\n3 third\n4 fourth\n5 fifth\n", spaces: true},
	"printf 'one\\n\\ntwo\\n' | nl -b a":                     {output: "1 one\n2\n3 two\n", spaces: true},
	"printf 'one\\n\\ntwo\\n' | nl -b t":                     {output: "1 one\n\n2 two\n", spaces: true},
	"nl -s'. ' notes.txt":                                    {output: "1. old old\n2. second\n3. third\n4. fourth\n5. fifth\n", spaces: true},
	"nl -w 2 notes.txt":                                      {output: "1 old old\n2 second\n3 third\n4 fourth\n5 fifth\n", spaces: true},
	"nl -n rz -w 3 notes.txt":                                {output: "001\told old\n002\tsecond\n003\tthird\n004\tfourth\n005\tfifth\n"},
	"echo 'the quick brown fox' | fold -w 10":                {output: "the quick \nbrown fox\n"},
	"echo 'the quick brown fox' | fold -s -w 12":             {output: "the quick \nbrown fox\n"},
	"echo 'the quick brown fox' | fold -w 12":                {output: "the quick br\nown fox\n"},
	"echo 'the quick brown fox' | fold -s -w 10 | grep -c .": {output: "2\n"},
	"printf 'abcdefgh\\n' | fold -b -w 4":                    {output: "abcd\nefgh\n"},
	// env, printenv, id, whoami, uname, hostname, df, yes, sleep, expand, unexpand:
	// results that depend on the machine or user are reduced to a stable shape.
	"env FOO=bar sh -c 'echo \"$FOO\"'":                                      {output: "bar\n"},
	"env -i FOO=bar env":                                                     {output: "FOO=bar\n"},
	"HOME=/x env -u HOME sh -c 'echo ${HOME:-unset}'":                        {output: "unset\n"},
	"env | grep -c '^PATH='":                                                 {output: "1\n"},
	"env -0 | tr '\\0' '\\n' | grep -c '^PATH='":                             {output: "1\n"},
	"FOO=bar printenv FOO":                                                   {output: "bar\n"},
	"printenv NOPE || echo unset":                                            {output: "unset\n"},
	"printenv | grep -c '^PATH='":                                            {output: "1\n"},
	"FOO=bar sh -c 'echo \"value is $(printenv FOO)\"'":                      {output: "value is bar\n"},
	"id -u | grep -c '^[0-9][0-9]*$'":                                        {output: "1\n"},
	"test \"$(id -un)\" = \"$(whoami)\" && echo same":                        {output: "same\n"},
	"id -g | grep -c '^[0-9][0-9]*$'":                                        {output: "1\n"},
	"id -Gn | grep -c .":                                                     {output: "1\n"},
	"{ test \"$(id -u)\" -eq 0 && echo root || echo not root; } | grep -c .": {output: "1\n"},
	"whoami | grep -c .":                                                     {output: "1\n"},
	"test \"$(whoami)\" = \"$(id -un)\" && echo same":                        {output: "same\n"},
	"echo \"/home/$(whoami)\" | grep -c '^/home/.'":                          {output: "1\n"},
	"test \"$(uname)\" = \"$(uname -s)\" && echo same":                       {output: "same\n"},
	"uname -s | grep -c .":                                                   {output: "1\n"},
	"uname -r | grep -c .":                                                   {output: "1\n"},
	"uname -m | grep -c .":                                                   {output: "1\n"},
	"case \"$(uname -s)\" in Darwin) echo mac;; Linux) echo linux;; *) echo other;; esac | grep -c .": {output: "1\n"},
	"test -n \"$(hostname)\" && echo set":                                       {output: "set\n"},
	"hostname -s | grep -c .":                                                   {output: "1\n"},
	"echo \"backup-$(hostname -s).tar\" | grep -c '^backup-.'":                  {output: "1\n"},
	"df -k . | head -n 1 | cut -d' ' -f1":                                       {output: "Filesystem\n"},
	"df -Pk . | head -n 1 | cut -d' ' -f1":                                      {output: "Filesystem\n"},
	"df -h . | head -n 1 | cut -d' ' -f1":                                       {output: "Filesystem\n"},
	"df -Pk . | grep -c .":                                                      {output: "2\n"},
	"df -Pk . | awk 'NR == 2 {print $5}' | grep -c '%$'":                        {output: "1\n"},
	"yes | head -n 3":                                                           {output: "y\ny\ny\n"},
	"yes ok | head -n 2":                                                        {output: "ok\nok\n"},
	"yes data | head -n 1000 | grep -c data":                                    {output: "1000\n"},
	"yes | (read a; read b; echo \"$a$b\")":                                     {output: "yy\n"},
	"sleep 1 && echo done":                                                      {output: "done\n"},
	"sleep 0 && echo ok":                                                        {output: "ok\n"},
	"for i in 1 2; do echo \"step $i\"; sleep 1; done":                          {output: "step 1\nstep 2\n"},
	"(sleep 1; echo later) & wait":                                              {output: "later\n"},
	"sleep 0.1 && echo done":                                                    {output: "done\n"},
	"printf 'a\\tb\\n' | expand":                                                {output: "a       b\n"},
	"printf 'a\\tb\\n' | expand -t 4":                                           {output: "a   b\n"},
	"printf 'x\\ty\\n' > in.txt && expand in.txt > out.txt && cat out.txt":      {output: "x       y\n"},
	"printf '\\ta\\tb\\n' | expand -i -t 4 | od -c | grep -c '\\\\t'":           {output: "1\n"},
	"printf '        a\\n' | unexpand | od -c | grep -c '\\\\t'":                {output: "1\n"},
	"printf 'a       b\\n' | unexpand -a | od -c | grep -c '\\\\t'":             {output: "1\n"},
	"printf 'a   b\\n' | unexpand -t 4 | od -c | grep -c '\\\\t'":               {output: "1\n"},
	"printf 'a\\tb\\n' | expand | od -c | grep -q '\\\\t' || echo no tabs left": {output: "no tabs left\n"},
	"printf 'a       b\\n' | unexpand | od -c | grep -q '\\\\t' || echo no tab": {output: "no tab\n"},
	"yes one two | head -n 2":                                                   {output: "one two\none two\n"},
	// checksum and encoding tools: digests of fixed input are the same everywhere.
	"printf abc | md5sum":                                         {output: "900150983cd24fb0d6963f7d28e17f72  -\n"},
	"printf abc | md5sum | cut -c1-8":                             {output: "90015098\n"},
	"printf abc > f.txt && md5sum f.txt":                          {output: "900150983cd24fb0d6963f7d28e17f72  f.txt\n"},
	"printf abc > f.txt && md5sum f.txt > SUMS && md5sum -c SUMS": {output: "f.txt: OK\n"},
	"printf abc > f.txt && md5sum f.txt > SUMS && printf x >> f.txt && md5sum -c SUMS 2> /dev/null || echo changed": {output: "f.txt: FAILED\nchanged\n"},
	"printf abc | sha1sum":                                          {output: "a9993e364706816aba3e25717850c26c9cd0d89d  -\n"},
	"printf abc | sha1sum | cut -c1-8":                              {output: "a9993e36\n"},
	"printf abc > f.txt && sha1sum f.txt":                           {output: "a9993e364706816aba3e25717850c26c9cd0d89d  f.txt\n"},
	"printf abc > f.txt && sha1sum f.txt > SUMS && sha1sum -c SUMS": {output: "f.txt: OK\n"},
	"printf abc > f.txt && sha1sum f.txt > SUMS && printf x >> f.txt && sha1sum -c SUMS 2> /dev/null || echo changed": {output: "f.txt: FAILED\nchanged\n"},
	"printf abc | sha224sum":                                            {output: "23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7  -\n"},
	"printf abc | sha224sum | cut -c1-8":                                {output: "23097d22\n"},
	"printf abc > f.txt && sha224sum f.txt":                             {output: "23097d223405d8228642a477bda255b32aadbce4bda0b3f7e36c9da7  f.txt\n"},
	"printf abc > f.txt && sha224sum f.txt > SUMS && sha224sum -c SUMS": {output: "f.txt: OK\n"},
	"printf abc > f.txt && sha224sum f.txt > SUMS && printf x >> f.txt && sha224sum -c SUMS 2> /dev/null || echo changed": {output: "f.txt: FAILED\nchanged\n"},
	"printf abc | sha256sum":                                            {output: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad  -\n"},
	"printf abc | sha256sum | cut -c1-8":                                {output: "ba7816bf\n"},
	"printf abc > f.txt && sha256sum f.txt":                             {output: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad  f.txt\n"},
	"printf abc > f.txt && sha256sum f.txt > SUMS && sha256sum -c SUMS": {output: "f.txt: OK\n"},
	"printf abc > f.txt && sha256sum f.txt > SUMS && printf x >> f.txt && sha256sum -c SUMS 2> /dev/null || echo changed": {output: "f.txt: FAILED\nchanged\n"},
	"printf abc | sha384sum":                                            {output: "cb00753f45a35e8bb5a03d699ac65007272c32ab0eded1631a8b605a43ff5bed8086072ba1e7cc2358baeca134c825a7  -\n"},
	"printf abc | sha384sum | cut -c1-8":                                {output: "cb00753f\n"},
	"printf abc > f.txt && sha384sum f.txt":                             {output: "cb00753f45a35e8bb5a03d699ac65007272c32ab0eded1631a8b605a43ff5bed8086072ba1e7cc2358baeca134c825a7  f.txt\n"},
	"printf abc > f.txt && sha384sum f.txt > SUMS && sha384sum -c SUMS": {output: "f.txt: OK\n"},
	"printf abc > f.txt && sha384sum f.txt > SUMS && printf x >> f.txt && sha384sum -c SUMS 2> /dev/null || echo changed": {output: "f.txt: FAILED\nchanged\n"},
	"printf abc | sha512sum":                                            {output: "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f  -\n"},
	"printf abc | sha512sum | cut -c1-8":                                {output: "ddaf35a1\n"},
	"printf abc > f.txt && sha512sum f.txt":                             {output: "ddaf35a193617abacc417349ae20413112e6fa4e89a97ea20a9eeee64b55d39a2192992a274fc1a836ba3c23a3feebbd454d4423643ce80e2a9ac94fa54ca49f  f.txt\n"},
	"printf abc > f.txt && sha512sum f.txt > SUMS && sha512sum -c SUMS": {output: "f.txt: OK\n"},
	"printf abc > f.txt && sha512sum f.txt > SUMS && printf x >> f.txt && sha512sum -c SUMS 2> /dev/null || echo changed": {output: "f.txt: FAILED\nchanged\n"},
	"printf abc | cksum":                {output: "1219131554 3\n"},
	"printf abc > f.txt && cksum f.txt": {output: "1219131554 3 f.txt\n"},
	"printf abc > a.txt && printf abc > b.txt && cksum a.txt b.txt | cut -d' ' -f1,2 | uniq | grep -c .": {output: "1\n"},
	"printf abc | cksum | cut -d' ' -f1; printf abd | cksum | cut -d' ' -f1":                             {output: "1219131554\n2137327320\n"},
	"printf abc | cksum -a sha256 | cut -c1-30":                                                          {output: "SHA256 (-) = ba7816bf8f01cfea4\n"},
	"printf hello | base64":       {output: "aGVsbG8=\n"},
	"printf aGVsbG8= | base64 -d": {output: "hello"},
	"printf hello > in.txt && base64 < in.txt > in.b64 && base64 -d < in.b64 > out.txt && cat out.txt": {output: "hello"},
	"printf 'a\\nb\\n' | base64 | base64 -d":                                                           {output: "a\nb\n"},
	"printf '%0100d' 0 | base64 -w 0 | grep -c .":                                                      {output: "1\n"},
	"printf '%0100d' 0 | base64 -b 76 | grep -c .":                                                     {output: "2\n"},
	"printf abc | b2sum | cut -c1-8":                                                                   {output: "ba80a53f\n"},
	"printf abc > f.txt && b2sum f.txt > SUMS && b2sum -c SUMS":                                        {output: "f.txt: OK\n"},
	"printf abc | b2sum -l 256 | cut -c1-8":                                                            {output: "bddd813c\n"},
	"printf hello | base32":                                                                            {output: "NBSWY3DP\n"},
	"printf NBSWY3DP | base32 -d":                                                                      {output: "hello"},
	"printf '%0100d' 0 | base32 -w 0 | grep -c .":                                                      {output: "1\n"},
	"printf hello | basenc --base64":                                                                   {output: "aGVsbG8=\n"},
	"printf hello | basenc --base16":                                                                   {output: "68656C6C6F\n"},
	"printf '\\377\\376' | basenc --base64url":                                                         {output: "__4=\n"},
	"printf aGVsbG8= | basenc --base64 -d":                                                             {output: "hello"},
	// tac, shuf, od, split, csplit, fmt, factor, expr, nproc, timeout, truncate:
	// random or machine-dependent output is reduced to a stable shape.
	"tac notes.txt":                                                              {output: "fifth\nfourth\nthird\nsecond\nold old\n"},
	"tac notes.txt | head -n 2":                                                  {output: "fifth\nfourth\n"},
	"printf 'a\\nb\\nc\\n' | tac":                                                {output: "c\nb\na\n"},
	"tac notes.txt | nl | head -n 1":                                             {output: "1 fifth\n", spaces: true},
	"seq 5 | shuf | sort -n | paste -sd, -":                                      {output: "1,2,3,4,5\n"},
	"shuf -e a b c | sort | paste -sd, -":                                        {output: "a,b,c\n"},
	"seq 6 > n6.txt && shuf -n 3 n6.txt | grep -c .":                             {output: "3\n"},
	"printf 'only\\n' | shuf -n 1":                                               {output: "only\n"},
	"shuf -i 1-5 | sort -n | paste -sd, -":                                       {output: "1,2,3,4,5\n"},
	"shuf -r -n 5 -e x | paste -sd, -":                                           {output: "x,x,x,x,x\n"},
	"printf abc | od -An -c":                                                     {output: "a b c\n", spaces: true},
	"printf abc | od -An -tx1":                                                   {output: "61 62 63\n", spaces: true},
	"printf abc | od -An -td1":                                                   {output: "97 98 99\n", spaces: true},
	"printf 'A\\n' | od -c | head -n 1":                                          {output: "0000000 A \\n\n", spaces: true},
	"printf 'a\\tb\\n' | od -An -c":                                              {output: "a \\t b \\n\n", spaces: true},
	"printf abc | od -An -tx1 -N 2":                                              {output: "61 62\n", spaces: true},
	"seq 6 > n6.txt && split -l 2 n6.txt part_ && ls part_*":                     {output: "part_aa\npart_ab\npart_ac\n"},
	"seq 6 > n6.txt && split -l 2 n6.txt part_ && cat part_aa":                   {output: "1\n2\n"},
	"seq 6 > n6.txt && split -l 2 n6.txt part_ && cat part_* | paste -sd, -":     {output: "1,2,3,4,5,6\n"},
	"printf abcdefgh | split -b 4 - chunk_ && cat chunk_aa":                      {output: "abcd"},
	"seq 6 > n6.txt && split -d -l 2 n6.txt p_ && ls p_*":                        {output: "p_00\np_01\np_02\n"},
	"seq 6 > n6.txt && split -n 3 n6.txt part_ && ls part_* | grep -c .":         {output: "3\n"},
	"seq 6 > n6.txt && csplit -s n6.txt 4 && paste -sd, xx00 && paste -sd, xx01": {output: "1,2,3\n4,5,6\n"},
	"seq 6 > n6.txt && csplit -s -f part_ n6.txt 4 && ls part_*":                 {output: "part_00\npart_01\n"},
	"printf 'a\\nb\\nc\\nb\\n' > words.txt && csplit -s words.txt '/c/' && paste -sd, xx00 && paste -sd, xx01": {output: "a,b\nc,b\n"},
	"printf 'a\\nb\\na\\nb\\n' > ab.txt && csplit -s ab.txt '/a/' '{*}' && ls xx* | grep -c .":                 {output: "3\n"},
	"echo 'the quick brown fox jumps' | fmt -w 15":                                                             {output: "the quick brown\nfox jumps\n"},
	"printf 'one\\ntwo\\nthree\\n' | fmt":                                                                      {output: "one two three\n"},
	"printf 'one\\ntwo\\n\\nthree\\nfour\\n' | fmt":                                                            {output: "one two\n\nthree four\n"},
	"factor 12":                        {output: "12: 2 2 3\n"},
	"factor 12 17":                     {output: "12: 2 2 3\n17: 17\n"},
	"echo 30 | factor":                 {output: "30: 2 3 5\n"},
	"factor 17 | awk '{print NF - 1}'": {output: "1\n"},
	"factor -h 12":                     {output: "12: 2^2 3\n"},
	"expr 2 + 3":                       {output: "5\n"},
	"expr 7 / 2":                       {output: "3\n"},
	"expr 3 \\* 4":                     {output: "12\n"},
	"expr 3 \\< 5":                     {output: "1\n"},
	"expr abc : 'a.*'":                 {output: "3\n"},
	"expr 1 - 1 || echo zero":          {output: "0\nzero\n"},
	"expr length hello":                {output: "5\n"},
	"nproc | grep -c '^[0-9][0-9]*$'":  {output: "1\n"},
	"jobs=$(nproc); test \"$jobs\" -ge 1 && echo ok":      {output: "ok\n"},
	"test \"$(nproc --all)\" -ge \"$(nproc)\" && echo ok": {output: "ok\n"},
	"timeout 5 echo hi":                                                             {output: "hi\n"},
	"timeout 1 sleep 5 || echo stopped":                                             {output: "stopped\n"},
	"timeout 5 sh -c 'exit 3'; echo $?":                                             {output: "3\n"},
	"{ timeout -s KILL 1 sleep 5; } 2> /dev/null || echo killed":                    {output: "killed\n"},
	"timeout 1 sleep 5; echo $?":                                                    {output: "124\n"},
	"truncate -s 10 f.txt && ls -l f.txt | awk '{print $5}'":                        {output: "10\n"},
	"truncate -s 0 notes.txt && ls -l notes.txt | awk '{print $5}'":                 {output: "0\n"},
	"truncate -s 5 new.bin && od -An -tx1 new.bin":                                  {output: "00 00 00 00 00\n", spaces: true},
	"truncate -s 3 notes.txt && cat notes.txt":                                      {output: "old"},
	"truncate -s 5 f.txt && truncate -s +5 f.txt && ls -l f.txt | awk '{print $5}'": {output: "10\n"},
	// dd, install, link, unlink, mkfifo, nohup, nice, sync, sum, tsort, pr, pathchk.
	"printf abcdef | dd bs=3 count=1 2>/dev/null":                              {output: "abc"},
	"printf abcdef | dd bs=1 skip=2 count=3 2>/dev/null":                       {output: "cde"},
	"dd if=notes.txt of=out.txt 2>/dev/null && head -n 1 out.txt":              {output: "old old\n"},
	"dd if=notes.txt of=out.txt bs=4096 2>/dev/null && head -n 1 out.txt":      {output: "old old\n"},
	"printf abc | dd conv=ucase status=none":                                   {output: "ABC"},
	"install -m 644 notes.txt copy.txt && ls -l copy.txt | cut -c1-10":         {output: "-rw-r--r--\n"},
	"install -m 755 notes.txt run.sh && ls -l run.sh | cut -c1-10":             {output: "-rwxr-xr-x\n"},
	"install -d x/y && ls -d x/y":                                              {output: "x/y\n"},
	"install -d -m 700 pv && ls -ld pv | cut -c1-10":                           {output: "drwx------\n"},
	"install -D notes.txt q/r/n.txt && ls q/r":                                 {output: "n.txt\n"},
	"link notes.txt hard.txt && head -n 1 hard.txt":                            {output: "old old\n"},
	"link notes.txt hard.txt && echo extra >> hard.txt && tail -n 1 notes.txt": {output: "extra\n"},
	"link notes.txt hard.txt && test notes.txt -ef hard.txt && echo same file": {output: "same file\n"},
	"unlink app.log && { ls app.log 2> /dev/null || echo gone; }":              {output: "gone\n"},
	"ln notes.txt hard.txt && unlink hard.txt && head -n 1 notes.txt":          {output: "old old\n"},
	"ln -s notes.txt link && unlink link && head -n 1 notes.txt":               {output: "old old\n"},
	"mkfifo p && test -p p && echo fifo":                                       {output: "fifo\n"},
	"mkfifo p && ls -l p | cut -c1":                                            {output: "p\n"},
	"mkfifo p && (echo hi > p &) && cat p":                                     {output: "hi\n"},
	"nohup echo hi":                                                            {output: "hi\n"},
	"nohup sh -c 'echo hi' > out.log 2>&1 && cat out.log":                      {output: "hi\n"},
	"nohup sleep 1 > /dev/null 2>&1 & wait; echo started":                      {output: "started\n"},
	"nice -n 5 echo hi":                                      {output: "hi\n"},
	"nice -n 10 sh -c 'echo working'":                        {output: "working\n"},
	"nice -n 5 sh -c 'exit 3'; echo $?":                      {output: "3\n"},
	"a=$(nice); b=$(nice -n 5 nice); echo $((b - a))":        {output: "5\n"},
	"sync && echo done":                                      {output: "done\n"},
	"printf data > f.txt && sync && cat f.txt":               {output: "data"},
	"sync; echo flushed":                                     {output: "flushed\n"},
	"printf data > f.txt && sync f.txt && echo ok":           {output: "ok\n"},
	"printf abc | sum":                                       {output: "16556 1\n", spaces: true},
	"printf abc > f.txt && sum f.txt | awk '{print $1, $2}'": {output: "16556 1\n"},
	"printf abc > a.txt && printf abc > b.txt && test \"$(sum a.txt | cut -d' ' -f1)\" = \"$(sum b.txt | cut -d' ' -f1)\" && echo same": {output: "same\n"},
	"printf 'a b\\nb c\\n' | tsort":                                      {output: "a\nb\nc\n"},
	"printf 'a b\\na c\\n' | tsort | head -n 1":                          {output: "a\n"},
	"printf 'a a\\nb c\\n' | tsort | sort | paste -sd, -":                {output: "a,b,c\n"},
	"printf 'a\\nb\\n' | pr -t":                                          {output: "a\nb\n"},
	"seq 4 | pr -t -2":                                                   {output: "1 3\n2 4\n", spaces: true},
	"printf 'a\\nb\\n' > f1 && printf '1\\n2\\n' > f2 && pr -t -m f1 f2": {output: "a 1\nb 2\n", spaces: true},
	"pathchk -p good/name && echo ok":                                    {output: "ok\n"},
	"pathchk -p 'a b' 2> /dev/null || echo not portable":                 {output: "not portable\n"},
	"pathchk -p /usr/bin/x && echo ok":                                   {output: "ok\n"},
	// groups, logname, tty, users, who, uptime, hostid, pinky, chown, chgrp, shred, dir, vdir,
	// dircolors, mknod: results that depend on the machine or login are reduced to a shape.
	"groups | grep -c .":                                 {output: "1\n"},
	"groups | wc -l | tr -d ' '":                         {output: "1\n"},
	"{ logname || echo none; } 2> /dev/null | grep -c .": {output: "1\n"},
	"{ test \"$(logname 2> /dev/null)\" = \"$(id -un)\" && echo same || echo different; } | grep -c .": {output: "1\n"},
	"tty < /dev/null; echo $?":                                                   {output: "not a tty\n1\n"},
	"tty < /dev/null || echo not a terminal":                                     {output: "not a tty\nnot a terminal\n"},
	"tty -s < /dev/null || echo not a terminal":                                  {output: "not a terminal\n"},
	"if tty > /dev/null < /dev/null; then echo interactive; else echo batch; fi": {output: "batch\n"},
	"users > /dev/null && echo ok":                                               {output: "ok\n"},
	"users | cat > /dev/null && echo ok":                                         {output: "ok\n"},
	"who > /dev/null && echo ok":                                                 {output: "ok\n"},
	"who -H > /dev/null && echo ok":                                              {output: "ok\n"},
	"who am i > /dev/null && echo ok":                                            {output: "ok\n"},
	"uptime | grep -c up":                                                        {output: "1\n"},
	"uptime | grep -c load":                                                      {output: "1\n"},
	"hostid | grep -c '^[0-9a-f]\\{8\\}$'":                                       {output: "1\n"},
	"pinky > /dev/null && echo ok":                                               {output: "ok\n"},
	"pinky -s > /dev/null && echo ok":                                            {output: "ok\n"},
	"chown \"$(id -u)\" notes.txt && echo ok":                                    {output: "ok\n"},
	"chown \"$(id -u):$(id -g)\" notes.txt && echo ok":                           {output: "ok\n"},
	"chown \"$(id -u)\" notes.txt && test \"$(ls -l notes.txt | awk '{print $3}')\" = \"$(id -un)\" && echo yes":       {output: "yes\n"},
	"chown -R \"$(id -u)\" projects logs && echo ok":                                                                   {output: "ok\n"},
	"chgrp \"$(id -g)\" notes.txt && echo ok":                                                                          {output: "ok\n"},
	"chgrp -R \"$(id -g)\" projects logs && echo ok":                                                                   {output: "ok\n"},
	"chgrp \"$(id -g)\" notes.txt && test \"$(ls -l notes.txt | awk '{print $4}')\" = \"$(id -gn)\" && echo yes":       {output: "yes\n"},
	"printf secret > s.txt && cp s.txt orig.txt && shred s.txt && { cmp -s s.txt orig.txt || echo overwritten; }":      {output: "overwritten\n"},
	"printf secret > s.txt && shred -u s.txt && { ls s.txt 2> /dev/null || echo removed; }":                            {output: "removed\n"},
	"printf secret > s.txt && cp s.txt orig.txt && shred -n 1 s.txt && { cmp -s s.txt orig.txt || echo overwritten; }": {output: "overwritten\n"},
	"mkdir d && touch d/a d/b && dir d":                                                                                {output: "a  b\n"},
	"mkdir d && touch 'd/a b' && dir d":                                                                                {output: "a\\ b\n"},
	"mkdir d && touch d/a d/b && vdir d | grep -c .":                                                                   {output: "3\n"},
	"mkdir d && touch d/a d/b && vdir d | tail -n +2 | awk '{print $NF}'":                                              {output: "a\nb\n"},
	"dircolors -b | grep -c '^LS_COLORS='":                                                                             {output: "1\n"},
	"dircolors -c | grep -c '^setenv LS_COLORS'":                                                                       {output: "1\n"},
	"dircolors -p | grep -c '^DIR'":                                                                                    {output: "1\n"},
	"mknod p p && test -p p && echo fifo":                                                                              {output: "fifo\n"},
	"mknod p p && ls -l p | cut -c1":                                                                                   {output: "p\n"},
	// echo, printf, pwd, true, false, test, kill, arch: the Bash builtin and the executables
	// must give the same result for the Common examples.
	"echo hello":                       {output: "hello\n"},
	"echo one two":                     {output: "one two\n"},
	"echo \"a   b\"":                   {output: "a   b\n"},
	"echo \"$((1 + 2))\"":              {output: "3\n"},
	"echo hi > out.txt && cat out.txt": {output: "hi\n"},
	"echo":                             {output: "\n"},
	"echo -n hi; echo":                 {output: "hi\n"},
	"echo -e 'a\\tb' | od -An -c":      {output: "a \\t b \\n\n", spaces: true},
	"printf '%s-%s\\n' a b":            {output: "a-b\n"},
	"printf '%5.2f|\\n' 3.14159":       {output: " 3.14|\n"},
	"printf '%d\\n' 0x1f":              {output: "31\n"},
	"printf '%s\\n' a b c":             {output: "a\nb\nc\n"},
	"printf 'a\\tb\\n'":                {output: "a\tb\n"},
	"printf '%b' 'a\\tb\\n'":           {output: "a\tb\n"},
	"printf -v greeting '%s' hi && echo \"$greeting\"":     {output: "hi\n"},
	"printf '%q\\n' 'a b'":                                 {output: "a\\ b\n"},
	"pwd | grep -c '^/'":                                   {output: "1\n"},
	"cd logs && pwd | grep -c 'logs$'":                     {output: "1\n"},
	"ln -s logs lnk && cd lnk && pwd -L | grep -c 'lnk$'":  {output: "1\n"},
	"ln -s logs lnk && cd lnk && pwd -P | grep -c 'logs$'": {output: "1\n"},
	"true; echo $?":                                        {output: "0\n"},
	"true && echo yes":                                     {output: "yes\n"},
	"false || true; echo continued":                        {output: "continued\n"},
	"n=0; while true; do n=$((n + 1)); test \"$n\" -ge 3 && break; done; echo \"$n\"": {output: "3\n"},
	"if true; then true; fi; echo ok":                                                 {output: "ok\n"},
	"false; echo $?":                                                                  {output: "1\n"},
	"false || echo failed":                                                            {output: "failed\n"},
	"if false; then echo a; else echo b; fi":                                          {output: "b\n"},
	"false && echo skipped; echo end":                                                 {output: "end\n"},
	"test -d logs && echo dir":                                                        {output: "dir\n"},
	"test -f notes.txt && echo file":                                                  {output: "file\n"},
	"test -e missing.txt || echo missing":                                             {output: "missing\n"},
	"test -s notes.txt && echo nonempty":                                              {output: "nonempty\n"},
	"test -z \"\" && echo empty":                                                      {output: "empty\n"},
	"test \"a\" = \"a\" && echo equal":                                                {output: "equal\n"},
	"test 3 -lt 5 && echo yes":                                                        {output: "yes\n"},
	"test ! -e missing.txt && echo not there":                                         {output: "not there\n"},
	"[ -n \"x\" ] && echo set":                                                        {output: "set\n"},
	"{ sleep 5 & kill $!; wait $!; } 2> /dev/null; echo done":                         {output: "done\n"},
	"kill -0 $$ && echo alive":                                                        {output: "alive\n"},
	"{ sleep 5 & kill -s TERM $!; wait $!; } 2> /dev/null; echo done":                 {output: "done\n"},
	"{ sleep 5 & kill -9 $!; wait $!; } 2> /dev/null; echo done":                      {output: "done\n"},
	"arch | grep -c '^[a-z0-9_]*$'":                                                   {output: "1\n"},
	"test -n \"$(arch)\" && echo named":                                               {output: "named\n"},
	// numfmt, ptx, stdbuf, coreutils run; stty, chroot, chcon, runcon need a terminal,
	// administrator rights or SELinux, so they are reviewed only.
	"numfmt --to=si 1500":                                      {output: "1.5k\n"},
	"numfmt --to=iec 1048576":                                  {output: "1.0M\n"},
	"numfmt --from=si 2K":                                      {output: "2000\n"},
	"echo 1234567 | numfmt --to=si":                            {output: "1.3M\n"},
	"numfmt --to=iec --suffix=B 2048":                          {output: "2.0KB\n"},
	"printf '1000 apples\\n' | numfmt --field 1 --to=si":       {output: "1.0k apples\n"},
	"numfmt --format %.2f 3.14159":                             {output: "3.15\n"},
	"printf 'the quick fox\\n' | ptx -w 40":                    {output: "the quick fox\nthe quick fox\nthe quick fox\n", spaces: true},
	"printf 'the quick fox\\n' | ptx -w 40 | grep -c .":        {output: "3\n"},
	"stdbuf -o0 echo hi":                                       {output: "hi\n"},
	"printf 'a\\nb\\n' | stdbuf -oL cat":                       {output: "a\nb\n"},
	"echo hi | stdbuf -i0 cat":                                 {output: "hi\n"},
	"stty -a":                                                  {notRun: true},
	"stty size":                                                {notRun: true},
	"stty -echo; read -r secret; stty echo":                    {notRun: true},
	"stty sane":                                                {notRun: true},
	"chroot /srv/root /bin/sh":                                 {notRun: true},
	"chroot /srv/root /bin/ls /":                               {notRun: true},
	"chroot --userspec=nobody:nogroup /srv/root /bin/id":       {notRun: true},
	"chcon -t httpd_sys_content_t /srv/www/index.html":         {notRun: true},
	"chcon --reference=/srv/www/index.html /srv/www/new.html":  {notRun: true},
	"chcon -R -t httpd_sys_content_t /srv/www":                 {notRun: true},
	"runcon -t httpd_t /usr/sbin/httpd":                        {notRun: true},
	"runcon system_u:system_r:httpd_t:s0 /usr/sbin/httpd":      {notRun: true},
	"runcon -t user_t id -Z":                                   {notRun: true},
	"coreutils --coreutils-prog=echo hello":                    {output: "hello\n"},
	"coreutils --coreutils-prog=basename logs/server.log .log": {output: "server\n"},
	"coreutils --coreutils-prog=true; echo $?":                 {output: "0\n"},
	// diff, cmp, diff3, sdiff, patch, zcat (and the gunzip example on the gzip page).
	"printf 'a\\nb\\n' > f1 && printf 'a\\nc\\n' > f2 && { diff -u f1 f2 | tail -n +3; }":                                                                                                              {output: "@@ -1,2 +1,2 @@\n a\n-b\n+c\n"},
	"printf 'a\\nb\\n' > f1 && printf 'a\\nc\\n' > f2 && diff -q f1 f2; echo \"status $?\"":                                                                                                            {output: "Files f1 and f2 differ\nstatus 1\n"},
	"printf 'a\\n' > f1 && cp f1 f2 && diff f1 f2 && echo same":                                                                                                                                        {output: "same\n"},
	"printf 'a\\n' > f1 && printf 'b\\n' > f2 && { diff -q f1 f2 > /dev/null || echo differ; }":                                                                                                        {output: "differ\n"},
	"printf 'x  y\\n' > w1 && printf 'x y\\n' > w2 && diff -w w1 w2 && echo same-ignoring-space":                                                                                                       {output: "same-ignoring-space\n"},
	"mkdir d1 d2 && printf 'x\\n' > d1/f && printf 'x\\n' > d2/f && diff -r d1 d2 && echo same":                                                                                                        {output: "same\n"},
	"printf 'a\\nb\\n' > f1 && printf 'a\\nc\\n' > f2 && diff f1 f2; test $? -eq 1 && echo differ":                                                                                                     {output: "2c2\n< b\n---\n> c\ndiffer\n"},
	"printf 'a\\nb\\n' > f1 && cp f1 f2 && cmp -s f1 f2 && echo same":                                                                                                                                  {output: "same\n"},
	"printf 'a\\nb\\n' > f1 && printf 'a\\nc\\n' > f2 && { cmp -s f1 f2 || echo differ; }":                                                                                                             {output: "differ\n"},
	"printf 'a\\n' > f1 && cmp -s f1 f1 && echo same":                                                                                                                                                  {output: "same\n"},
	"printf 'a\\n' > f1 && printf 'a\\nb\\n' > f2 && { cmp -s f1 f2 || echo differ; }":                                                                                                                 {output: "differ\n"},
	"printf 'a\\nb\\n' > f1 && printf 'a\\nc\\n' > f2 && cmp f1 f2; test $? -eq 1 && echo differ":                                                                                                      {output: "f1 f2 differ: char 3, line 2\ndiffer\n"},
	"printf 'a\\nb\\n' > mine && printf 'a\\nb\\n' > base && printf 'a\\nb\\nd\\n' > yours && diff3 mine base yours":                                                                                   {output: "====3\n1:2a\n2:2a\n3:3c\n  d\n"},
	"printf 'A\\nb\\n' > mine && printf 'a\\nb\\n' > base && printf 'a\\nb\\nd\\n' > yours && diff3 -m mine base yours":                                                                                {output: "A\nb\nd\n"},
	"printf 'a\\nb\\n' > f1 && printf 'a\\nc\\n' > f2 && sdiff f1 f2; test $? -eq 1 && echo differ":                                                                                                    {output: "a a\nb | c\ndiffer\n", spaces: true},
	"printf 'a\\nb\\n' > f1 && printf 'a\\nc\\n' > f2 && sdiff -w 20 f1 f2; test $? -eq 1 && echo differ":                                                                                              {output: "a a\nb | c\ndiffer\n", spaces: true},
	"printf 'a\\nb\\n' > f1 && printf 'a\\nc\\n' > f2 && sdiff -s f1 f2; test $? -eq 1 && echo differ":                                                                                                 {output: "b | c\ndiffer\n", spaces: true},
	"printf 'a\\nb\\n' > old.txt && printf 'a\\nc\\n' > new.txt && { diff -u old.txt new.txt > change.patch; patch old.txt change.patch > /dev/null && cat old.txt; }":                                 {output: "a\nc\n"},
	"printf 'a\\nb\\n' > old.txt && printf 'a\\nc\\n' > new.txt && { diff -u old.txt new.txt > change.patch; cp new.txt target.txt; patch -R target.txt change.patch > /dev/null && cat target.txt; }": {output: "a\nb\n"},
	"printf 'a\\nb\\n' > old.txt && printf 'a\\nc\\n' > new.txt && { diff -u old.txt new.txt > change.patch; patch --dry-run old.txt change.patch > /dev/null && cat old.txt; }":                       {output: "a\nb\n"},
	"printf 'x\\n' > f && printf 'y\\n' > g && { diff -u --label a/f --label b/f f g > d.patch; patch -p1 < d.patch > /dev/null && cat f; }":                                                           {output: "y\n"},
	"printf data | gzip > f.gz && zcat < f.gz":                                                                                                                                                         {output: "data"},
	"printf 'a\\nb\\n' | gzip > f.gz && zcat < f.gz | wc -l | tr -d ' '":                                                                                                                               {output: "2\n"},
	"printf data | gzip > f.gz && gunzip -c f.gz":                                                                                                                                                      {output: "data"},
	"printf data | gzip > f.gz && zcat f.gz":                                                                                                                                                           {output: "data"},
	"printf data | gzip > f.gz && gzcat f.gz":                                                                                                                                                          {output: "data"},
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
	lines := strings.Split(strings.TrimRight(s, "\n"), "\n")
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
