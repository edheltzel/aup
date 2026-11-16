package config

import "testing"

func TestSuggestionsCatchTypos(t *testing.T) {
	src := []byte("[[step]]\nlabel = \"a\"\nbinary = \"b\"\nversion = \"--version\"\ntimout = \"5s\"\ncolour = \"#fff\"\n")
	got := Suggestions("steps.toml", src, nil)
	want := map[string]bool{`timout`: false, `colour`: false}
	for _, s := range got {
		for k := range want {
			if contains(s, k) {
				want[k] = true
			}
		}
	}
	for k, seen := range want {
		if !seen {
			t.Errorf("%s not suggested in %v", k, got)
		}
	}
	if far := Suggestions("steps.toml", []byte("[[step]]\nlabel=\"a\"\nbinary=\"b\"\nversion=\"-v\"\ncompletely = 1\n"), nil); len(far) != 0 {
		t.Fatalf("a distant key was suggested: %v", far)
	}
}

func TestLegacyFilenameIsFlagged(t *testing.T) {
	got := Suggestions("/home/x/.config/aup/harnesses.toml", []byte("[[step]]\n"), nil)
	if len(got) != 1 || !contains(got[0], "steps.toml") {
		t.Fatalf("got %v", got)
	}
}

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
