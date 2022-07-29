package example

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
	"regexp"
)

func ExampleReplaceFile() {
	var filepath = "./example.txt"
	bfile.SetContents(filepath, "example content")

	fmt.Println(bfile.GetContents(filepath))
	bfile.ReplaceFile("example content", "replace content", filepath)
	fmt.Println(bfile.GetContents(filepath))
	// output:
	// example content
	// replace content
}

func ExampleReplaceFileFunc() {

	filepath := bfile.Pwd() + bfile.Separator + "test.txt"
	bfile.SetContents(filepath, "123456")
	fmt.Println(bfile.GetContents(filepath))
	bfile.ReplaceFileFunc(func(path, content string) string {
		reg, _ := regexp.Compile(`\d[2-4]`)
		return reg.ReplaceAllString(content, "[num]")
	}, filepath)
	fmt.Println(bfile.GetContents(filepath))
	bfile.Remove(filepath)

	// output:
	// 123456
	// [num][num]56
}

func ExampleReplaceDir() {
	subDir := bfile.Pwd() + bfile.Separator + "sub"
	file := subDir + bfile.Separator + "test.txt"
	file2 := subDir + bfile.Separator + "test2.txt"
	bfile.SetContents(file, "test 123")
	bfile.SetContents(file2, "test 456")
	fmt.Printf("file:%s\n", bfile.GetContents(file))
	fmt.Printf("file2:%s\n", bfile.GetContents(file2))
	bfile.ReplaceDir("test", "file", subDir, "test*", true)
	fmt.Printf("file:%s\n", bfile.GetContents(file))
	fmt.Printf("file2:%s\n", bfile.GetContents(file2))

	bfile.Remove(subDir)

	// output:
	// file:test 123
	// file2:test 456
	// file:file 123
	// file2:file 456
}

func ExampleReplaceDirFunc() {
	subDir := bfile.Pwd() + bfile.Separator + "sub"
	file := subDir + bfile.Separator + "test.txt"
	file2 := subDir + bfile.Separator + "test2.txt"
	bfile.SetContents(file, "test 123")
	bfile.SetContents(file2, "test 456")
	fmt.Printf("file:%s\n", bfile.GetContents(file))
	fmt.Printf("file2:%s\n", bfile.GetContents(file2))
	bfile.ReplaceDirFunc(func(path, content string) string {
		reg, _ := regexp.Compile(`\d[0-9]`)
		return reg.ReplaceAllString(content, "num")
	}, subDir, "test*", true)
	fmt.Printf("file:%s\n", bfile.GetContents(file))
	fmt.Printf("file2:%s\n", bfile.GetContents(file2))

	bfile.Remove(subDir)

	// output:
	// file:test 123
	// file2:test 456
	// file:test num3
	// file2:test num6
}
