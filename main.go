package main

import "fmt"

type Tmp struct {
}

type rob struct {
	nodes map[string]Tmp
}

var std = &rob{
	nodes: make(map[string]Tmp),
}

func main() {
	t := Tmp{}
	fmt.Printf("%p\n", &t)
	std.nodes["123"] = t
	v := std.nodes["123"]
	fmt.Printf("%p\n", &v)

}
