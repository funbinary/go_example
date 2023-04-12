// github.com/bigwhite/experiments/tree/master/slog-examples/demo2/main.go
package main

import (
	"fmt"
)

func main() {
	var ch chan *int
	go func() {
		<-ch
	}()
	select {
	case ch <- nil:
		fmt.Println("it's time")
	}

}
