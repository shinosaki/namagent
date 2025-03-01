package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func Escape(input string, replace string) string {
	pattern := `[\\:;*"<>|&#!?%@+=^~]`
	re := regexp.MustCompile(pattern)
	return re.ReplaceAllString(input, replace)
}

func MkDir(filename string) (dir string, err error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return "", err
	}
	dir = filepath.Dir(path)
	return dir, os.MkdirAll(dir, os.ModePerm)
}

func UniqueFilename(path string) string {
	check := func(path string) bool {
		if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
			return true
		}
		return false
	}

	if check(path) {
		return path
	}

	ext := filepath.Ext(path)
	base := strings.TrimSuffix(path, ext)

	i := 1
	for {
		new := fmt.Sprintf("%s_%d%s", base, i, ext)
		if check(new) {
			return new
		}
		i++
	}
}
