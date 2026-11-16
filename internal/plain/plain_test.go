package plain

import (
	"bytes"
	"strings"
	"testing"
)

func TestStripperHandlesSplitAndCharsetEscapes(t *testing.T) {
	var log bytes.Buffer
	s := &Stripper{W: &log}
	for _, chunk := range []string{"a\x1b[3", "1mred\x1b[0m", "\x1b(Bplain\x1b[m", "\x1b]0;ti", "tle\x1b", "\\end"} {
		if n, err := s.Write([]byte(chunk)); err != nil || n != len(chunk) {
			t.Fatalf("write %q: n=%d err=%v", chunk, n, err)
		}
	}
	if got := log.String(); got != "aredplainend" {
		t.Fatalf("log=%q", got)
	}
	if strings.Contains(log.String(), "\x1b") {
		t.Fatalf("log has escape codes: %q", log.String())
	}
}
