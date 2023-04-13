package main

import (
	"encoding/asn1"
	"fmt"
)

func main() {
	mdata, err := asn1.Marshal(0x3e)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(mdata)
}
