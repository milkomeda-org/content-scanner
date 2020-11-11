// Created by vinson on 2020/11/4.

package utils

import (
	"io/ioutil"
	"os"
	"path"
	"sort"
	"text/template"
)

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func IsDir(path string) (bool, *os.FileInfo) {
	s, err := os.Stat(path)
	if err != nil {
		return false, nil
	}
	return s.IsDir(), &s
}

func PathAssemble(p, f string) string {
	if path.Base(p) == "/" || path.Base(p) == "\\" {
		return p + f
	}
	return p + "/" + f
}

// 读取文件字符串
func ReadAll(filePth string) (string, error) {
	f, err := os.Open(filePth)
	if err != nil {
		return "", err
	}
	readAll, err := ioutil.ReadAll(f)
	if nil == err {
		return string(readAll), err
	}
	return "", err
}

func GetTemplate(path string) (*template.Template, error) {
	files := []string{path}
	return template.ParseFiles(files...)
}

func ReadDirOrderByModify(dirname string) ([]os.FileInfo, error) {
	fs, err := ioutil.ReadDir(dirname)
	if err != nil {
		return fs, err
	}
	return sortByTime(fs), nil
}

func sortByTime(pl []os.FileInfo) []os.FileInfo {
	sort.Slice(pl, func(i, j int) bool {
		flag := false
		if pl[i].ModTime().After(pl[j].ModTime()) {
			flag = true
		} else if pl[i].ModTime().Equal(pl[j].ModTime()) {
			if pl[i].Name() < pl[j].Name() {
				flag = true
			}
		}
		return flag
	})
	return pl
}
