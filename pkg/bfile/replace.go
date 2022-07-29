package bfile

import "strings"

// ReplaceFile
//  @Description: 替换文件search的内容为replace
//  @param search
//  @param replace
//  @param path
//  @return error
//
func ReplaceFile(search, replace, path string) error {
	return SetContents(path, strings.Replace(GetContents(path), search, replace, -1))
}

// ReplaceFileFunc
//  @Description: 使用自定义的函数替换文件内容
//
//  @param f
//
//  @param path
//
//  @return error
//
func ReplaceFileFunc(f func(path, content string) string, path string) error {
	data := GetContents(path)
	result := f(path, data)
	if len(data) != len(result) && data != result {
		return SetContents(path, result)
	}
	return nil
}

// ReplaceDir
//  @Description: 扫描文件路径，将符合条件的文件的指定内容替换为新内容
//
//  @param search 匹配的内容
//
//  @param replace 需要替代的字符串
//
//  @param path 文件路径
//
//  @param pattern 匹配的文件名
//
//  @param recursive 是否递归
//
//  @return error
//
func ReplaceDir(search, replace, path, pattern string, recursive ...bool) error {
	files, err := ScanDirFile(path, pattern, recursive...)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = ReplaceFile(search, replace, file); err != nil {
			return err
		}
	}
	return err
}

// ReplaceDirFunc
//  @Description: 扫描指定目录，使用自定义函数替换符合条件的文件的指定内容为新内容
//
//  @param f
//
//  @param path 文件路径
//
//  @param pattern 匹配的文件名
//
//  @param recursive 是否递归
//
//  @return error
//
func ReplaceDirFunc(f func(path, content string) string, path, pattern string, recursive ...bool) error {
	files, err := ScanDirFile(path, pattern, recursive...)
	if err != nil {
		return err
	}
	for _, file := range files {
		if err = ReplaceFileFunc(f, file); err != nil {
			return err
		}
	}
	return err
}
