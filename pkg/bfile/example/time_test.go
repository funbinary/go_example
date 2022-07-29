package example

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
)

func ExampleMTime() {
	fmt.Println(bfile.MTime("D:\\workspace\\libpcap-1.10.1-ubuntu-x64.tar.gz"))

	// output:
	// 2021-10-19 10:13:42 +0800 CST
}

func ExampleMTimeStamp() {
	fmt.Println(bfile.MTimestamp("D:\\workspace\\libpcap-1.10.1-ubuntu-x64.tar.gz"))

	// output:
	// 1634609622
}

func ExampleMTimestampMilli() {
	fmt.Println(bfile.MTimestampMilli("D:\\workspace\\libpcap-1.10.1-ubuntu-x64.tar.gz"))

	// output:
	// 1634609622000
}
