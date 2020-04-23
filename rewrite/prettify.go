package rewrite

import (
	"go/format"
	"io/ioutil"

	"golang.org/x/tools/imports"
)

const DEFAULT_FILE_MODE = 0644

func WriteCodeFile(fileName string, sourceCode []byte) ([]byte, error) {
	err := ioutil.WriteFile(fileName, sourceCode, DEFAULT_FILE_MODE)
	return sourceCode, err
}

func WriteGolangFile(fileName string, sourceCode []byte) ([]byte, error) {
	// Formart/Prettify the code 格式化代码
	srcCode, err := format.Source(sourceCode)
	if err != nil {
		return sourceCode, err
	}
	if _, err = WriteCodeFile(fileName, srcCode); err != nil {
		return srcCode, err
	}
	// Split the imports in two groups: go standard and the third parts 分组排序引用包
	var dstCode []byte
	dstCode, err = imports.Process(fileName, srcCode, nil)
	if err != nil {
		return srcCode, err
	}
	return WriteCodeFile(fileName, dstCode)
}
