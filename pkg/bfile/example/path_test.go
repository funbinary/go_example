package example_test

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
)

func ExampleJoin() {
	fmt.Println(bfile.Join("D:\\workspace\\kds\\base_package\\bfile\\example\\", "log"))

	// Output:
	// D:\workspace\kds\base_package\bfile\example\log
}

func ExampleBasename() {
	fmt.Println(bfile.Basename("/var/www/webrtc.js"))
	fmt.Println(bfile.Basename("/var/www/"))
	fmt.Println(bfile.Basename("webrtc.js"))

	// Output:
	// webrtc.js
	// www
	// webrtc.js
}

func ExampleName() {
	fmt.Println(bfile.Name("/var/www/webrtc.js"))
	fmt.Println(bfile.Name("/var/www/"))
	fmt.Println(bfile.Name("webrtc.js"))

	// Output:
	// webrtc
	// www
	// webrtc
}

func ExampleExt() {
	fmt.Println(bfile.Ext("/var/www/webrtc.js"))
	fmt.Println(bfile.Ext("/var/www/"))
	fmt.Println(bfile.Ext("webrtc.js"))
	fmt.Println(bfile.ExtName("/var/www/webrtc.js"))
	fmt.Println(bfile.ExtName("/var/www/"))
	fmt.Println(bfile.ExtName("webrtc.js"))
	fmt.Println(bfile.ExtName(".js"))
	fmt.Println(bfile.ExtName("js"))

	// Output:
	// .js
	//
	//.js
	//js
	//
	//js

}

func ExampleAbs() {
	fmt.Println(bfile.Abs("/kds/data/kds"))
	fmt.Println(bfile.Abs("/kds/data/kds.exe"))
	fmt.Println(bfile.Abs("./main.go"))

	// Output:
	// D:\kds\data\kds
	// D:\kds\data\kds.exe
	// D:\workspace\kds\base_package\bfile\example\main.go
}

func ExampleRealPath() {

	fmt.Println(bfile.RealPath("D:\\workspace\\kds\\base_package\\go.mod"))

	// Output:
	// D:\workspace\kds\base_package\go.mod
}

func ExampleDir() {

	fmt.Println(bfile.Dir("/foo/bar/baz.js"))
	fmt.Println(bfile.Dir("/foo/bar/baz"))
	fmt.Println(bfile.Dir("/foo/bar/baz/"))
	fmt.Println(bfile.Dir("/dirty//path///"))
	fmt.Println(bfile.Dir("dev.txt"))
	fmt.Println(bfile.Dir("../data/todo.txt"))
	fmt.Println(bfile.Dir(".."))
	fmt.Println(bfile.Dir("."))
	fmt.Println(bfile.Dir("/"))
	fmt.Println(bfile.Dir(""))
	fmt.Println("中文:", bfile.Dir("/中文测试"))

	// Output:
	// \foo\bar
	// \foo\bar
	// \foo\bar\baz
	// \dirty\path
	// .
	// ..\data
	// .
	// D:\workspace\kds\base_package\bfile
	// \
	// .
}
