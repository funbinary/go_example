package example_test

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
	"log"
	"os"
)

func ExampleMkdir() {
	//例子中使用相对路径，但是实际推荐使用绝对路径
	err := bfile.Mkdir("./test/asd")
	fmt.Printf("%+v", err)
}

func ExampleCreate() {
	_, err := bfile.Create("./test/fff/a")
	fmt.Println(err)
	_, err = bfile.Create("./test/fff/a/2.txt")
	fmt.Println(err)

	// Output:
	// <nil>
	// os.Create failed for name "./test/fff/a/2.txt"
}

func ExampleRemove() {
	bfile.Mkdir("./test/a")

	fmt.Println(bfile.Remove("./test/a"))
	// Output:
	// <nil>

}

func ExampleExists() {
	fmt.Println(bfile.Exists("D:/workspace"))
	fmt.Println(bfile.Exists("D:\\workspace"))
	fmt.Println(bfile.Exists("D:\\workspace\\kds\\base_package\\"))
	fmt.Println(bfile.Exists("D:\\workspace\\kds\\base_package\\noexist"))
	fmt.Println(bfile.Exists("D:/workspace/kds/base_package"))
	fmt.Println(bfile.Exists("D:/workspace/kds/base_package/go.mod"))
	fmt.Println(bfile.Exists("./go.mod"))

	// Output:
	// true
	// true
	// true
	// false
	// true
	// true
	// false

}

func ExampleOpen() {
	f, err := bfile.Open("./noexit")

	fmt.Printf("open file fail:%v", err) //使用%v可以打印错误。使用%+v可以打印堆栈信息
	if err == nil {
		f.Close()
	}

	// output:
	// open file fail:os.Open failed for name "./noexit"
}

func ExampleOpenFile() {
	f, err := bfile.OpenFile("./noexit", bfile.O_RDONLY, 0777)

	fmt.Printf("open file fail:%v\n", err) //使用%v可以打印错误。使用%+v可以打印堆栈信息
	if err == nil {
		f.Close()
	}

	f, err = bfile.OpenFile("./noexit", bfile.O_CREATE, 0777)
	fmt.Printf("open file fail:%v", err) //使用%v可以打印错误。使用%+v可以打印堆栈信息

	if err == nil {
		f.Close()
	}

	// output:
	// open file fail:os.OpenFile failed with name "./noexit", flag "0", perm "511"
	// open file fail:<nil>
}

func ExampleStat() {
	info, err := bfile.Stat("bfile_test.go")
	if err != nil {
		fmt.Printf("%v", err) //使用%v可以打印错误。使用%+v可以打印堆栈信息
	}
	fmt.Println(info.ModTime())
	fmt.Println(info.IsDir())
	fmt.Println(info.Size())
	fmt.Println(info.Sys())
	fmt.Println(info.Name())
	fmt.Println(info.Mode())

	info2, err := bfile.Stat("./noexist")

	if err != nil {
		fmt.Printf("%v", err) //使用%v可以打印错误。使用%+v可以打印堆栈信息
	} else {
		fmt.Println(info2.ModTime())
		fmt.Println(info2.IsDir())
		fmt.Println(info2.Size())
		fmt.Println(info2.Sys())
		fmt.Println(info2.Name())
		fmt.Println(info2.Mode())
	}

	// output:
	// 2022-02-28 11:21:20.1893404 +0800 CST
	// false
	// 2560
	// &{32 {1612290097 30935943} {1095523356 30944338} {1095523356 30944338} 0 2560}
	// bfile_test.go
	// -rw-rw-rw-
	// os.Stat failed for file "./noexist"

}

func ExampleMove() {
	//rename for file
	src := "test.txt"
	dst := "test2.txt"
	f, _ := bfile.Create(src)
	f.Close()

	err := bfile.Move(src, dst)
	fmt.Printf("%v\n", err)
	os.Remove(dst)

	//move for dir
	srcDir := "test"
	dstDir := "test2"
	bfile.Mkdir(srcDir)
	err = bfile.Move(srcDir, dstDir)
	fmt.Printf("%v\n", err)
	os.Remove(dstDir)

	// output:
	// <nil>
	// <nil>
}

func ExampleRename() {
	//rename for file
	src := "test.txt"
	dst := "test2.txt"
	f, _ := bfile.Create(src)
	f.Close()

	err := bfile.Rename(src, dst)
	fmt.Printf("%v\n", err)
	os.Remove(dst)

	//rename for dir
	srcDir := "test"
	bfile.Mkdir(srcDir)
	dstDir := "test2"
	err = bfile.Rename(srcDir, dstDir)
	fmt.Printf("%v\n", err)
	os.Remove(dstDir)

	// output:
	// <nil>
	// <nil>
}

func ExampleIsEmpty() {

	// check for empty file
	target := "1.txt"
	f, err := os.Create(target)
	if err != nil {
		log.Panicln(err)
	}
	f.Close()
	fmt.Println(bfile.IsEmpty(target))
	// check for no empty file
	bfile.SetContents(target, "123")
	fmt.Println(bfile.IsEmpty(target))
	os.Remove(target)

	// check for dir
	dir := "test"
	sub := "test/sub"

	bfile.Mkdir(dir)
	bfile.Mkdir(sub)
	fmt.Println(bfile.IsEmpty(dir))
	fmt.Println(bfile.IsEmpty(sub))
	os.Remove(dir)
	os.Remove(sub)

	// output:
	// true
	// false
	// false
	// true
}

func ExampleDirNames() {
	sub := "test/sub"
	noexist := "noexist"
	bfile.Mkdir(sub)
	bfile.Mkdir(sub + "/ss")
	bfile.Mkdir("test/sub2")
	bfile.SetContents("test/1.txt", "test")
	fmt.Println(bfile.DirSubs("test"))
	fmt.Println(bfile.DirSubs(noexist))
	fmt.Println(bfile.DirSubs("./bfile_test.go"))
	os.Remove("test")
	// output:
	// [1.txt sub sub2] <nil>
	// [] os.Open failed for name "noexist"
	// [] Read dir files failed from path "./bfile_test.go"
}

func ExampleGlob() {
	path := bfile.Pwd() + bfile.Separator + "bfile_test.go"
	match, _ := bfile.Glob(path, true)
	fmt.Println(match)
	match, _ = bfile.Glob(path, false)
	fmt.Println(match)
	bfile.Chdir("/")
	match, _ = bfile.Glob(path, false)
	fmt.Println(match)
	// output:
	// [bfile_test.go]
	// [D:\workspace\kds\base_package\bfile\example\bfile_test.go]
	// [D:\workspace\kds\base_package\bfile\example\bfile_test.go]
}

func ExampleIsReadable() {

	path := bfile.Pwd() + bfile.Separator + "bfile_test.go"
	fmt.Println(bfile.IsReadable(path))

	// output:
	// true
}

func ExampleIsWriteable() {

	path := bfile.Pwd() + bfile.Separator + "bfile_test.go"
	fmt.Println(bfile.IsWritable(path))

	// output:
	// true
}
