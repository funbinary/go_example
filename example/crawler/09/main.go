package main

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"time"
)

// xpath
func main() {
	// 1. 创建谷歌浏览器实例
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	// 2. 设置context超时时间
	ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	var example string
	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://pkg.go.dev/time`),
		chromedp.WaitVisible(`body > footer`),
		chromedp.Click(`#example-After`, chromedp.NodeVisible),
		chromedp.Value(`#example-After textarea`, &example),
	)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Go's time.After example:\\n%s", example)

}
