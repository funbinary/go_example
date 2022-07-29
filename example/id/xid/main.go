package main

import (
	"fmt"

	"github.com/rs/xid"
)

func main() {
	guid := xid.New()

	fmt.Println(guid.Value())
	fmt.Println(guid.Pid())
	fmt.Println(guid.String())
	fmt.Println(guid.Machine())

}
