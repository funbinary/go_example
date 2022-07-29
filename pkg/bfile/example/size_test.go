package example

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
)

func ExampleSize() {
	fmt.Println(bfile.Size("D:\\workspace"))
	fmt.Println(bfile.Size("D:\\workspace\\libpcap-1.10.1-ubuntu-x64.tar.gz"))

	// output:
	// 4096
	// 905197
}

func ExampleSizeFormat() {
	fmt.Println(bfile.SizeFormat("D:\\workspace"))
	fmt.Println(bfile.SizeFormat("D:\\workspace\\libpcap-1.10.1-ubuntu-x64.tar.gz"))

	// output:
	// 4.00K
	// 883.98K
}

func ExampleSizeReadableSize() {
	fmt.Println(bfile.ReadableSize("D:\\workspace"))
	fmt.Println(bfile.ReadableSize("D:\\workspace\\libpcap-1.10.1-ubuntu-x64.tar.gz"))

	// output:
	// 4.00K
	// 883.98K
}

func ExampleStrToSize() {
	fmt.Println(bfile.StrToSize("800K"))

	// output:
	// 819200
}

func ExampleFormatSize() {
	fmt.Println(bfile.FormatSize(500000000))

	// output:
	// 476.84M
}
