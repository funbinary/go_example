package example

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
	"strings"
)

func ExampleScanDir() {

	fmt.Println(bfile.ScanDir(bfile.Pwd(), "bfile_test.go", false))
	fmt.Println(bfile.ScanDir(bfile.Pwd(), "noexist", true))
	fmt.Println(bfile.ScanDir(bfile.Pwd(), "*", true))

	// output:
	// [D:\workspace\kds\base_package\bfile\example\bfile_test.go] <nil>
	// [] <nil>
	// [D:\workspace\kds\base_package\bfile\example\bfile_test.go D:\workspace\kds\base_package\bfile\example\cmd_test.go D:\workspace\kds\base_package\bfile\example\content_test.go D:\workspace\kds\base_package\bfile\example\copy_test.go D:\workspace\kds\base_package\bfile\example\path_test.go D:\workspace\kds\base_package\bfile\example\replace_test.go D:\workspace\kds\base_package\bfile\example\scan_test.go D:\workspace\kds\base_package\bfile\example\self_test.go D:\workspace\kds\base_package\bfile\example\type_test.go] <nil>
}

func ExampleScanDirFunc() {

	fmt.Println(bfile.ScanDirFunc(bfile.Pwd(), "*", true, func(path string) string {
		newPath := strings.ReplaceAll(path, bfile.Pwd(), "")
		return newPath
	}))

	// output:
	// [D:\bfile_test.go D:\cmd_test.go D:\content_test.go D:\copy_test.go D:\path_test.go D:\replace_test.go D:\scan_test.go D:\self_test.go D:\type_test.go] <nil>

}

func ExampleScanDirFile() {

	fmt.Println(bfile.ScanDirFile(bfile.Pwd(), "*", true))

	// output:
	// [D:\workspace\kds\base_package\bfile\example\bfile_test.go D:\workspace\kds\base_package\bfile\example\cmd_test.go D:\workspace\kds\base_package\bfile\example\content_test.go D:\workspace\kds\base_package\bfile\example\copy_test.go D:\workspace\kds\base_package\bfile\example\path_test.go D:\workspace\kds\base_package\bfile\example\replace_test.go D:\workspace\kds\base_package\bfile\example\scan_test.go D:\workspace\kds\base_package\bfile\example\self_test.go D:\workspace\kds\base_package\bfile\example\type_test.go] <nil>

}

func ExampleScanDirFileFunc() {

	fmt.Println(bfile.ScanDirFileFunc(bfile.Pwd(), "*", true, func(path string) string {
		newPath := strings.ReplaceAll(path, bfile.Pwd(), "")
		return newPath
	}))

	// output:
	// [D:\bfile_test.go D:\cmd_test.go D:\content_test.go D:\copy_test.go D:\path_test.go D:\replace_test.go D:\scan_test.go D:\self_test.go D:\type_test.go] <nil>

}
