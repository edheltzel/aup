package config

import (
	"fmt"
	"strings"
)

var knownKeys = []string{"label", "binary", "installed", "version", "args", "env", "timeout", "color", "interactive"}

func Suggestions(path string, src []byte, bad []RowError) []string {
	var out []string
	if strings.HasSuffix(path, "harnesses.toml") {
		out = append(out, path+": rename this file to steps.toml. aup still reads the old name, but it is going away")
	}
	for _, k := range unknownKeys(src) {
		if s := nearest(k, knownKeys); s != "" {
			out = append(out, fmt.Sprintf("%q is not a setting. did you mean %q?", k, s))
		}
	}
	return out
}

func unknownKeys(src []byte) []string {
	var out []string
	in := false
	for _, line := range strings.Split(string(src), "\n") {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "[["):
			name := strings.Trim(line, "[]")
			in = name == "step" || name == "harness"
		case in && strings.Contains(line, "="):
			k := strings.TrimSpace(strings.SplitN(line, "=", 2)[0])
			if !known(k) {
				out = append(out, k)
			}
		}
	}
	return out
}

func known(k string) bool {
	for _, s := range knownKeys {
		if s == k {
			return true
		}
	}
	return false
}

func nearest(got string, known []string) string {
	got = strings.ToLower(got)
	best, dist := "", 2
	for _, k := range known {
		if d := editDistance(got, k); d > 0 && d < dist {
			best, dist = k, d
		}
	}
	return best
}

func editDistance(a, b string) int {
	if len(a) < len(b) {
		a, b = b, a
	}
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min(cur[j-1]+1, prev[j]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(b)]
}
