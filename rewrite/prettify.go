package rewrite

import (
	"go/format"
	"io/ioutil"

	"gitea.com/azhai/refactor/config"
	"gitea.com/azhai/refactor/utils"
	"golang.org/x/tools/imports"
)

// 格式化代码，如果出错返回原内容
func FormatGolangCode(src []byte) ([]byte, error) {
	_src, err := format.Source(src)
	if err == nil {
		src = _src
	}
	return src, err
}

func WriteCodeFile(fileName string, sourceCode []byte) ([]byte, error) {
	err := ioutil.WriteFile(fileName, sourceCode, config.DEFAULT_FILE_MODE)
	return sourceCode, err
}

func WriteGolangFile(fileName string, sourceCode []byte) ([]byte, error) {
	// Formart/Prettify the code 格式化代码
	srcCode, err := FormatGolangCode(sourceCode)
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

// 将包中的Go文件格式化，如果提供了pkgname则用作新包名
func RewritePackage(pkgpath, pkgname string) error {
	if pkgname != "" {
		// TODO: 替换包名
	}
	files, err := utils.FindFiles(pkgpath, ".go")
	if err != nil {
		return err
	}
	var content []byte
	for fileName := range files {
		content, err = ioutil.ReadFile(fileName)
		if err != nil {
			break
		}
		_, err = WriteGolangFile(fileName, content)
		if err != nil {
			break
		}
	}
	return err
}
