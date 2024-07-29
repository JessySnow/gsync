package helper

import (
	"os"
	"testing"
)

func TestMkdirIfAbsent(t *testing.T) {
	tmpDir := "/tmp/A/B/C"
	err := MkdirIfAbsent(tmpDir)
	if err != nil {
		t.Errorf("mkdir failed: %v", err)
		return
	}
}

func TestCreateFileIfAbsent(t *testing.T) {
	tmpFile := "/tmp/temp.txt"
	err := os.Remove(tmpFile)
	if err != nil && !os.IsNotExist(err) {
		t.Errorf("remove file failed: %v", err)
	}

	_, err = CreateFileIfAbsent(tmpFile)
	if err != nil {
		t.Errorf("create file failed: %v", err)
		return
	}

	_, err = os.Stat(tmpFile)
	if os.IsNotExist(err) {
		t.Errorf("file not exists: %v", tmpFile)
		return
	}
}
