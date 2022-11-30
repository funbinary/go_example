package main

/*
   gin框架实现文件下载功能
*/

import (
	"context"
	"fmt"
	"github.com/funbinary/go_example/pkg/errors"
)

//主函数
func main() {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := ctx.Err()
	if errors.Is(err, context.Canceled) {
		fmt.Println("cancel")
	}
}
