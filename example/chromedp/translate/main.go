// Command text is a chromedp example demonstrating how to extract text from a
// specific element.
package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"regexp"
)

var re = regexp.MustCompile(`^\d+\. .+$`) // 匹配一级标题

type SuricataDocInfo struct {
	Name string
	href string
}

func suricata() string {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck, //不检查默认浏览器
		//chromedp.Flag("headless", false),                 //开启图像界面,重点是开启这个
		chromedp.Flag("ignore-certificate-errors", true), //忽略错误
		chromedp.Flag("disable-web-security", true),      //禁用网络安全标志
		chromedp.NoFirstRun,                              //设置网站不是首次运行
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.164 Safari/537.36"), //设置UserAgent
	)
	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	//	defer cancel()
	ctx, cancel := chromedp.NewContext(
		allocCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	url := "https://suricata.readthedocs.io/en/suricata-6.0.9/"
	//err := chromedp.Run(ctx, chromedp.Navigate("https://suricata.readthedocs.io/en/suricata-6.0.9/"))

	var links []*cdp.Node

	var example string
	// run task list
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		//chromedp.Nodes(`.//a[@class="reference internal"]`, &links),
		chromedp.Nodes(`//*[@id="suricata-user-guide"]/.//a[@class="reference internal"]`, &links),
	)
	if err != nil {
		log.Fatal(err)
	}
	var infos []*SuricataDocInfo
	for _, link := range links {

		href := link.AttributeValue("href")
		if href == "" {
			continue
		}
		name := ""
		for _, sublink := range link.Children {
			if sublink.NodeName == "#text" {
				name = sublink.NodeValue
			}
		}
		if name == "" {
			continue
		}
		info := &SuricataDocInfo{
			Name: name,
			href: url + href,
		}
		infos = append(infos, info)
	}

	return example
}

func main() {
	fmt.Println(suricata())
}
