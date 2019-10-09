package mysqlclusterinit

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
)

// ReplaceFileContent 使用正则表达式查找模式，并且替换正则1号捕获分组为指定的内容
func ReplaceFileContent(filename, regexStr, repl string) error {
	conf, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("ReadFile %s error %w", filename, err)
	}

	fixed, err := ReplaceContent(string(conf), regexStr, repl)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, []byte(fixed), 0644)
}

// FileExists 检查文件是否存在，并且不是目录
func FileExists(filename string) error {
	if fi, err := os.Stat(filename); err != nil {
		return err
	} else if fi.IsDir() {
		return fmt.Errorf("file %s is a directory", filename)
	}

	return nil
}

// ReplaceContent 使用正则表达式查找模式，并且替换正则1号捕获分组为指定的内容
func ReplaceContent(str, regexStr, repl string) (string, error) {
	re, err := regexp.Compile(regexStr)
	if err != nil {
		return "", err
	}

	fixed := ""
	lastIndex := 0

	for _, v := range re.FindAllStringSubmatchIndex(str, -1) {
		if len(v) != 4 {
			return "", fmt.Errorf("regexp %s should have only one capturing group", regexStr)
		}

		fixed += str[lastIndex:v[2]] + repl
		lastIndex = v[3]
	}

	if lastIndex == 0 {
		return "", fmt.Errorf("regexp %s found non submatches", regexStr)
	}

	return fixed + str[lastIndex:], nil
}
