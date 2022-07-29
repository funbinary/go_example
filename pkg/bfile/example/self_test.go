package example_test

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
)

func ExampleSelf() {
	fmt.Println(bfile.SelfPath())
	fmt.Println(bfile.SelfName())
	fmt.Println(bfile.SelfDir())

	// Output:
	// D:\workspace\kds\base_package\bfile\example\go_build_beyondinfo_com_baselib_go_base_package_git_bfile_example.exe
	// go_build_beyondinfo_com_baselib_go_base_package_git_bfile_example.exe
	// D:\workspace\kds\base_package\bfile\example
}
