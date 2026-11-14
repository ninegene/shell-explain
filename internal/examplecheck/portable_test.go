package examplecheck

import "testing"

func TestPortableOutputVariations(t *testing.T) {
	cases := []struct {
		code  string
		valid []string
		bad   string
	}{
		{"echo 'the quick brown fox jumps' | fmt -w 15", []string{"the quick brown\nfox jumps\n", "the quick\nbrown fox\njumps\n"}, "the quick brown fox jumps\n"},
		{"numfmt --to=si 1500", []string{"1.5k\n", "1.5K\n"}, "1.5M\n"},
		{"printf '1000 apples\\n' | numfmt --field 1 --to=si", []string{"1.0k apples\n", "1.0K apples\n"}, "1.0K oranges\n"},
	}
	for _, tc := range cases {
		for _, output := range tc.valid {
			if err := Verify(tc.code, t.TempDir(), output); err != nil {
				t.Errorf("%s: rejected %q: %v", tc.code, output, err)
			}
		}
		if err := Verify(tc.code, t.TempDir(), tc.bad); err == nil {
			t.Errorf("%s: accepted invalid output %q", tc.code, tc.bad)
		}
	}
}
