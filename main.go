package main

import (
	"fmt"
)

func fff() <-chan struct{} {
	ch := make(chan struct{}, 100)
	ch <- struct{}{}
	fmt.Println(len(ch))

	return ch
}

func main() {
	ch := fff()
	fmt.Println(len(ch))

}
