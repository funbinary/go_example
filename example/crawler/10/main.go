package main

import (
	"fmt"
	"github.com/funbinary/go_example/example/crawler/10/collect"
	"github.com/funbinary/go_example/example/crawler/10/proxy"

	"time"
)

// xpath
func main() {
	// 1. 创建谷歌浏览器实例
	//ctx, cancel := chromedp.NewContext(context.Background())
	//defer cancel()
	//
	//// 2. 设置context超时时间
	//ctx, cancel = context.WithTimeout(ctx, 15*time.Second)
	//defer cancel()
	//
	//var example string
	//err := chromedp.Run(ctx,
	//	chromedp.Navigate(`https://pkg.go.dev/time`),
	//	chromedp.WaitVisible(`body > footer`),
	//	chromedp.Click(`#example-After`, chromedp.NodeVisible),
	//	chromedp.Value(`#example-After textarea`, &example),
	//)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("Go's time.After example:\\n%s", example)

	proxyURLs := []string{"http://127.0.0.1:10809", "http://127.0.0.1:10809"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		//logger.Error("RoundRobinProxySwitcher failed")
	}
	url := "https://google.com"
	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}

	body, err := f.Get(url)
	if err != nil {
		//logger.Error("read content failed",
		//	zap.Error(err),
		//)
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}
