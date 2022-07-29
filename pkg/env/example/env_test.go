package example

import (
	"fmt"
	"reflect"

	"github.com/funbinary/go_example/pkg/env"
)

func ExampleAll() {
	fmt.Println(env.All())

	// output:
	//
}

func ExampleMap() {
	for k, v := range env.Map() {
		fmt.Println(k, v)
	}

	// output:
	//
}

func ExampleGet() {
	fmt.Println(env.Get("CC"), reflect.TypeOf("CC"))

	fmt.Println(env.Get("NUMBER_OF_PROCESSORS"), reflect.TypeOf(env.Get("NUMBER_OF_PROCESSORS")))
	fmt.Println(env.Get("NOT", "all"))
	fmt.Println(env.Get("NOT_INT", 1))

	// output:
	//
}

func ExampleSet() {

	env.Set("Example", "test")
	fmt.Println(env.Get("Example"))
	fmt.Println(env.Get("EXAMPLE"))

	// output:
	//
}

func ExampleSetMap() {

	env.SetMap(map[string]string{
		"Example1": "test",
		"Example2": "test2",
	})

	fmt.Println(env.Get("Example1"))
	fmt.Println(env.Get("EXAMPLE2"))

	// output:
	//
}

func ExampleContains() {

	fmt.Println(env.Contains("CC"))
	fmt.Println(env.Contains("EXAMPLE2"))

	// output:
	//
}

func ExampleRemove() {
	fmt.Println(env.Get("MYSQL_IP"))

	fmt.Println(env.Remove("MYSQL_IP"))
	fmt.Println(env.Get("MYSQL_IP"))

	// output:
	//
}
