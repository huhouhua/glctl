package editor

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestEditor(t *testing.T) {
	edit := Editor{Args: []string{"cat"}}
	testStr := "test something\n"
	contents, path, err := edit.LaunchTempFile("", "someprefix", bytes.NewBufferString(testStr))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("no temp file: %s", path)
	}
	defer os.Remove(path)
	if disk, err := os.ReadFile(path); err != nil || !bytes.Equal(contents, disk) {
		t.Errorf("unexpected file on disk: %v %s", err, string(disk))
	}
	if !bytes.Equal(contents, []byte(testStr)) {
		t.Errorf("unexpected contents: %s", string(contents))
	}
	if !strings.Contains(path, "someprefix") {
		t.Errorf("path not expected: %s", path)
	}
}
