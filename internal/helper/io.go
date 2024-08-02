package helper

import (
	"fmt"
	"gync/internal/llog"
	"os"
)

// MkdirIfAbsent mkdir if release dir is absent, dir is the absolute dir path
func MkdirIfAbsent(dir string) error {
	if len(dir) == 0 {
		return fmt.Errorf("invalid dir path, which is empty string")
	}

	// create dir if absent
	if _, e := os.Stat(dir); os.IsNotExist(e) {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	} else {
		llog.Infof("dir %s already exists", dir)
	}

	return nil
}

// CreateFileIfAbsent create file if file is absent, file is the absolute file path
func CreateFileIfAbsent(file string) (f *os.File, err error) {
	if len(file) == 0 {
		return nil, fmt.Errorf("invalid file path, which is empty string")
	}

	if _, e := os.Stat(file); os.IsNotExist(e) {
		f, err := os.Create(file)
		if err != nil {
			return nil, err
		}
		return f, nil
	} else {
		llog.Infof("file %s already exists", file)
	}

	return os.Open(file)
}
