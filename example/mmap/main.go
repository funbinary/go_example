package main

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/mmap"
	expmmap "golang.org/x/exp/mmap"
	"os"
)

func Expmmap() {
	f, _ := os.OpenFile("tmp.txt", os.O_CREATE|os.O_RDWR, 0644)
	_, _ = f.WriteAt([]byte("abcdefg"), 0)
	defer f.Close()

	at, _ := expmmap.Open("tmp.txt")

	buff := make([]byte, 2)
	_, _ = at.ReadAt(buff, 4)
	_ = at.Close()
	fmt.Println(string(buff))
}

func Selfmmap() {

	f, err := os.OpenFile("notes.txt", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	_, err = f.WriteString("Hello, world")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	m, err := mmap.Map(f, mmap.RDWR, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(m))
	if err := m.Unmap(); err != nil {
		panic(err)
	}
}

func main() {
	Expmmap()

	Selfmmap()
}
