package config

import (
	"os"
	"path/filepath"
)

func load(filename string) (buf []byte, err error) {
	buf, err = os.ReadFile(filename)
	return
}

func save(filename string, buf []byte) (err error) {
	err = os.MkdirAll(filepath.Dir(filename), 0777)
	if err != nil {
		return
	}
	err = os.WriteFile(filename, buf, 0666)
	return
}
