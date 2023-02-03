package goutil

import (
	"fmt"
	"github.com/gookit/goutil/stdutil"
)

func ExampleBaseTypeVal() {
	if val, err := stdutil.BaseTypeVal(2); err == nil {
		fmt.Println(val)
	}
	fmt.Println(stdutil.GoVersion())

	// output:
	// 2
	// 1.19.4
}
