package example

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bfile"
	"os"
)

func ExampleContent() {
	var filepath = "./example.txt"
	//bfile.Create(filepath)
	bfile.SetContents(filepath, "123")
	fmt.Println(bfile.GetContents(filepath))
	bfile.SetContents(filepath, "456")
	fmt.Println(bfile.GetContents(filepath))
	bfile.AppendContents(filepath, "789")
	fmt.Println(bfile.GetContents(filepath))

	// Output:
	// 123
	// 456
	// 456789
}

func ExampleGetCharFromOffset() {
	var filepath = "./example.txt"

	bfile.SetContents(filepath, "123456789")
	f, _ := bfile.OpenFile(filepath, os.O_RDONLY, bfile.DefaultPermOpen)
	defer f.Close()

	index := bfile.GetCharFromOffset(f, '4', 0)
	fmt.Println(index)
	index = bfile.GetCharFromOffset(f, '4', 5)
	fmt.Println(index)
	index = bfile.GetCharFromOffset(f, '4', 100)
	fmt.Println(index)

	// Output:
	// 3
	// -1
	// -1

}

func ExampleGetCharOffsetFromByPath() {
	var filepath = "./example.txt"

	bfile.SetContents(filepath, "123456789")

	index := bfile.GetCharOffsetFromByPath(filepath, '4', 0)
	fmt.Println(index)
	index = bfile.GetCharOffsetFromByPath(filepath, '4', 5)
	fmt.Println(index)
	index = bfile.GetCharOffsetFromByPath(filepath, '4', 100)
	fmt.Println(index)

	// Output:
	// 3
	// -1
	// -1

}

func ExampleGetBytesByTwoOffsets() {
	var filepath = "./example.txt"

	bfile.SetContents(filepath, "123456789")

	f, _ := bfile.OpenFile(filepath, os.O_RDONLY, bfile.DefaultPermOpen)

	s := bfile.GetBytesByRange(f, 1, 5)
	fmt.Println(s)
	fmt.Println(string(s))

	// Output:
	// [50 51 52 53]
	// 2345

}

func ExampleGetBytesByTwoOffsetsByPath() {
	var filepath = "./example.txt"

	bfile.SetContents(filepath, "123456789")
	s := bfile.GetBytesByRangesByPath(filepath, 1, 5)
	fmt.Println(s)
	fmt.Println(string(s))

	// Output:
	// [50 51 52 53]
	// 2345
}

func ExampleGetBytesTilCharByPath() {
	var filepath = "./example.txt"

	bfile.SetContents(filepath, "123456789")

	chars, num := bfile.GetBytesTilCharByPath(filepath, 0, '4')
	fmt.Println(chars, num)
	fmt.Println(string(chars), num)

	chars, num = bfile.GetBytesTilCharByPath(filepath, 4, '4')
	fmt.Println(chars, num)
	fmt.Println(string(chars), num)
	chars, num = bfile.GetBytesTilCharByPath(filepath, -12, '4')
	fmt.Println(chars, num)
	fmt.Println(string(chars), num)

	// Output:
	// [49 50 51 52] 3
}

func ExampleReadLines() {

	var filepath = "./example.txt"

	bfile.SetContents(filepath, "12345\n6789")
	bfile.ReadLines(filepath, func(text string) error {
		fmt.Println(text)
		return nil
	})

	// OUTPUT:
	// 12345
	// 6789
}

func ExampleReadLinesBytes() {
	var filepath = "./example.txt"

	bfile.SetContents(filepath, "12345\n6789")
	bfile.ReadLinesBytes(filepath, func(bytes []byte) error {
		fmt.Println(bytes)
		return nil
	})

	// OUTPUT:
	// [49 50 51 52 53]
	// [54 55 56 57]
}

func ExampleTruncate() {
	var filepath = "./example.txt"

	bfile.SetContents(filepath, "abcdefghi")
	stat, _ := os.Stat(filepath)
	fmt.Println(stat.Size())
	bfile.Truncate(filepath, 1)
	stat, _ = os.Stat(filepath)
	fmt.Println(stat.Size())
	fmt.Println(bfile.GetContents(filepath))

	// Output:
	// 9
	// 1
	// a
}
