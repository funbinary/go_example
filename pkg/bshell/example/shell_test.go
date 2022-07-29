package example

import (
	"fmt"
	"github.com/funbinary/go_example/pkg/bshell"
	"log"
)

func ExampleShell() {
	result, err := bshell.ShellExec(`echo 1`)
	if err != nil {
		log.Panicln(err)
	}
	fmt.Println(result)

	// output:
	//
}
