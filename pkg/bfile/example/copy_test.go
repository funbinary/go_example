package example

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
)

func ExampleCopy() {
	bfile.Mkdir("test")
	bfile.Mkdir("test/a")
	bfile.Mkdir("test/b")
	bfile.SetContents("test/test.txt", "asd123")

	bfile.Copy("test/test.txt", "test/test2.txt")
	fmt.Println(bfile.GetContents("test/test2.txt"))

	bfile.Copy("test", "test2")
	sub, _ := bfile.DirSubs("test2")
	bfile.Remove("test")
	bfile.Remove("test2")
	fmt.Println(sub)

	// output:
	// asd123
	// [a b test.txt test2.txt]
}

func ExampleCopyFile() {

	bfile.Mkdir("test")
	bfile.Mkdir("test/a")
	bfile.Mkdir("test/b")
	bfile.SetContents("test/test.txt", "asd123")

	bfile.CopyFile("test/test.txt", "test/test2.txt")
	fmt.Println(bfile.GetContents("test/test2.txt"))

	bfile.CopyFile("test", "test2")
	sub, _ := bfile.DirSubs("test3")
	bfile.Remove("test")
	bfile.Remove("test2")
	fmt.Println(sub)

	// output:
	// asd123
	// []
}

func ExampleCopyDir() {

	bfile.Mkdir("test")
	bfile.Mkdir("test/a")
	bfile.Mkdir("test/b")
	bfile.SetContents("test/test.txt", "asd123")

	bfile.CopyDir("test/test.txt", "test/test2.txt")
	fmt.Println(bfile.GetContents("test/test2.txt"))

	bfile.CopyDir("test", "test2")
	sub, _ := bfile.DirSubs("test3")
	bfile.Remove("test")
	bfile.Remove("test2")
	fmt.Println(sub)

	// output:
	// asd123
	// [a b test.txt test2.txt]
}
