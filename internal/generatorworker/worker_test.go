package generatorworker

import (
	"encoding/base64"
	"strings"
	"testing"
)

func TestRunScriptFlagContract(t *testing.T) {
	want := []string{"flag{with,comma}", "flag{中文}", "flag{quote'and\"newline\n}"}
	outer, err := base64.StdEncoding.DecodeString(EncodeFlags(want))
	if err != nil {
		t.Fatal(err)
	}
	parts := strings.Split(string(outer), ",")
	if len(parts) != len(want) {
		t.Fatalf("flag boundaries lost: %d", len(parts))
	}
	for i, part := range parts {
		value, err := base64.StdEncoding.DecodeString(part)
		if err != nil || string(value) != want[i] {
			t.Fatalf("flag %d was changed", i)
		}
	}
}
