package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)


// 遍历目录下的文件
func FindFiles(dir, ext string) (map[string]os.FileInfo, error) {
	var result = make(map[string]os.FileInfo)
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return result, err
	}
	for _, file := range files {
		fname := file.Name()
		if ext != "" && !strings.HasSuffix(fname, ext) {
			continue
		}
		fname = filepath.Join(dir, fname)
		result[fname] = file
	}
	return result, nil
}
