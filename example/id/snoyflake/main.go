package main

import (
	"fmt"
	"log"

	"github.com/sony/sonyflake"
)

func main() {
	var sf *sonyflake.Sonyflake

	var st sonyflake.Settings
	sf = sonyflake.NewSonyflake(st)
	if sf == nil {
		log.Panicln("sonyflake not created")
	}
	id, err := sf.NextID()
	if err != nil {
		log.Panicln(err)
	}

	fmt.Println(id)
	fmt.Println(sonyflake.Decompose(id))

}
