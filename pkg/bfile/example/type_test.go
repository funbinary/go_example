package example

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
)

func ExampleIsFile() {
	var filepath = "./example.txt"
	fmt.Println(bfile.IsFile("./noexist.txt"))
	bfile.Create(filepath)
	fmt.Println(bfile.IsFile(filepath))

	fmt.Println(bfile.IsFile(bfile.SelfDir()))

	// output:
	// false
	// true
	// false
}

func ExampleIsDir() {
	var filepath = "./test"
	fmt.Println(bfile.IsDir("./noexist"))
	bfile.Mkdir(filepath)
	fmt.Println(bfile.IsDir(filepath))

	fmt.Println(bfile.IsDir(bfile.SelfDir()))

	// output:
	// false
	// true
	// true
}
