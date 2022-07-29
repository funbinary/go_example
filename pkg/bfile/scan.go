package bfile

import (
	"path/filepath"
	"sort"
	"strings"

	"github.com/funbinary/go_example/pkg/errors"
)

const (
	// Max recursive depth for directory scanning.
	maxScanDepth = 100000
)

// ScanDir
//  @Description: 扫描指定目录或文件，支持递归扫描
//  @param path
//  @param pattern 支持多个匹配，可以使用,分割，其匹配语法
//	term:
//		'*'         匹配0或多个非路径分隔符的字符
//		'?'         匹配1个非路径分割符的字符
//		'[' [ '^' ] { character-range } ']'
//		           字符组，必须非空
//		c          匹配字符，字符不能是*、?、\\、[
//		'\\' c     匹配字符，支持*、?、\\、[
//
//	 character-range:
//		c           匹配字符，字符不能是\\、-、]
//		'\\' c      匹配字符
//		lo '-' hi   匹配区间[lo,hi]的字符
//  匹配要求匹配整个name字符串，而不是它的一部分。
//  @param recursive true递归扫描，false不进行递归扫描
//  @return []string
//  @return error
//
func ScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	list, err := doScanDir(0, path, pattern, isRecursive, nil)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFunc
//  @Description: 扫描指定目录或文件，支持递归扫描
//  @param path
//  @param pattern 支持多个匹配，可以使用,分割，其匹配语法
//	term:
//		'*'         匹配0或多个非路径分隔符的字符
//		'?'         匹配1个非路径分割符的字符
//		'[' [ '^' ] { character-range } ']'
//		           字符组，必须非空
//		c          匹配字符，字符不能是*、?、\\、[
//		'\\' c     匹配字符，支持*、?、\\、[
//
//	character-range:
//		c           匹配字符，字符不能是\\、-、]
//		'\\' c      匹配字符
//		lo '-' hi   匹配区间[lo,hi]的字符
// 匹配要求匹配整个name字符串，而不是它的一部分。
//  @param recursive true递归扫描，false不进行递归扫描
//  @param handler 自定义方法，每次搜索到文件或者目录将会调用此方法，如果返回的字符串为空此路径将会被过滤掉
//  @return []string
//  @return error
//
func ScanDirFunc(path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	list, err := doScanDir(0, path, pattern, recursive, handler)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFile
//  @Description: 扫描指定目录的文件，支持递归扫描
//  @param path
//  @param pattern 支持多个匹配，可以使用,分割，其匹配语法
//	term:
//		'*'         匹配0或多个非路径分隔符的字符
//		'?'         匹配1个非路径分割符的字符
//		'[' [ '^' ] { character-range } ']'
//		           字符组，必须非空
//		c          匹配字符，字符不能是*、?、\\、[
//		'\\' c     匹配字符，支持*、?、\\、[
//
//	character-range:
//		c           匹配字符，字符不能是\\、-、]
//		'\\' c      匹配字符
//		lo '-' hi   匹配区间[lo,hi]的字符
//  匹配要求匹配整个name字符串，而不是它的一部分。
//  @param recursive true递归扫描，false不进行递归扫描
//  @param handler 自定义方法，每次搜索到文件或者目录将会调用此方法，如果返回的字符串为空此路径将会被过滤掉
//  @return []string
//  @return error
//
func ScanDirFile(path string, pattern string, recursive ...bool) ([]string, error) {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	list, err := doScanDir(0, path, pattern, isRecursive, func(path string) string {
		if IsDir(path) {
			return ""
		}
		return path
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// ScanDirFileFunc
//  @Description: 扫描指定目录的文件，支持递归扫描
//  @param path
//  @param pattern 支持多个匹配，可以使用,分割，其匹配语法
//	term:
//		'*'         匹配0或多个非路径分隔符的字符
//		'?'         匹配1个非路径分割符的字符
//		'[' [ '^' ] { character-range } ']'
//		           字符组，必须非空
//		c          匹配字符，字符不能是*、?、\\、[
//		'\\' c     匹配字符，支持*、?、\\、[
//
//	character-range:
//		c           匹配字符，字符不能是\\、-、]
//		'\\' c      匹配字符
//		lo '-' hi   匹配区间[lo,hi]的字符
//  匹配要求匹配整个name字符串，而不是它的一部分。
//  @param recursive true递归扫描，false不进行递归扫描
//  @return []string
//  @return error
//
func ScanDirFileFunc(path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	list, err := doScanDir(0, path, pattern, recursive, func(path string) string {
		if IsDir(path) {
			return ""
		}
		return handler(path)
	})
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// doScanDir is an internal method which scans directory and returns the absolute path
// list of files that are not sorted.
//
// The pattern parameter `pattern` supports multiple file name patterns, using the ','
// symbol to separate multiple patterns.
//
// The parameter `recursive` specifies whether scanning the `path` recursively, which
// means it scans its sub-files and appends the files path to result array if the sub-file
// is also a folder. It is false in default.
//
// The parameter `handler` specifies the callback function handling each sub-file path of
// the `path` and its sub-folders. It ignores the sub-file path if `handler` returns an empty
// string, or else it appends the sub-file path to result slice.
func doScanDir(depth int, path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	if depth >= maxScanDepth {
		return nil, errors.Errorf("directory scanning exceeds max recursive depth: %d", maxScanDepth)
	}
	var (
		list      []string
		file, err = Open(path)
	)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		err = errors.Wrapf(err, `read directory files failed from path "%s"`, path)
		return nil, err
	}
	var (
		filePath = ""

		patterns = strings.Split(pattern, ",")
	)
	for _, name := range names {
		filePath = path + Separator + name
		if IsDir(filePath) && recursive {
			array, _ := doScanDir(depth+1, filePath, pattern, true, handler)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		// Handler filtering.
		if handler != nil {
			filePath = handler(filePath)
			if filePath == "" {
				continue
			}
		}
		// If it meets pattern, then add it to the result list.
		for _, p := range patterns {
			if match, _ := filepath.Match(p, name); match {
				if filePath = Abs(filePath); filePath != "" {
					list = append(list, filePath)
				}
			}
		}
	}
	return list, nil
}
