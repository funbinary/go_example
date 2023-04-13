// github.com/bigwhite/experiments/tree/master/slog-examples/demo2/main.go
package main

<<<<<<< .mine
import (
	"encoding/asn1"
	"fmt"
)
=======
import (
	"fmt"
)

>>>>>>> .theirs

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
