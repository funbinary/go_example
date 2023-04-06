// Command text is a chromedp example demonstrating how to extract text from a
// specific element.
package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

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
	for _, link := range links {

		href := link.AttributeValue("href")
		if href != "" {
			fmt.Println(href)
		}

		for _, sublink := range link.Children {
			if sublink.NodeName == "#text" {
				fmt.Println(sublink.NodeValue)
			}
		}
	}
	return example
}

func travelSubtree(pageUrl, of string, opts ...chromedp.QueryOption) chromedp.Tasks {
	var nodes []*cdp.Node
	return chromedp.Tasks{
		chromedp.Navigate(pageUrl),
		chromedp.Nodes(of, &nodes, opts...),
		// ask chromedp to populate the subtree of a node
		chromedp.ActionFunc(func(c context.Context) error {
			// depth -1 for the entire subtree
			// do your best to limit the size of the subtree
			return dom.RequestChildNodes(nodes[0].NodeID).WithDepth(-1).Do(c)
		}),
		// wait a little while for dom.EventSetChildNodes to be fired and handled
		chromedp.Sleep(time.Second),
		chromedp.ActionFunc(func(c context.Context) error {
			printNodes(os.Stdout, nodes, "", "  ")
			return nil
		}),
	}
}

func printNodes(w io.Writer, nodes []*cdp.Node, padding, indent string) {
	// This will block until the chromedp listener closes the channel
	for _, node := range nodes {
		switch {
		case node.NodeName == "#text":
			fmt.Fprintf(w, "%s#text: %q\n", padding, node.NodeValue)
		default:
			fmt.Fprintf(w, "%s%s:\n", padding, strings.ToLower(node.NodeName))
			if n := len(node.Attributes); n > 0 {
				fmt.Fprintf(w, "%sattributes:\n", padding+indent)
				for i := 0; i < n; i += 2 {
					fmt.Fprintf(w, "%s%s: %q\n", padding+indent+indent, node.Attributes[i], node.Attributes[i+1])
				}
			}
		}
		if node.ChildNodeCount > 0 {
			fmt.Fprintf(w, "%schildren:\n", padding+indent)
			printNodes(w, node.Children, padding+indent+indent, indent)
		}
	}
}

func main() {
	fmt.Println(suricata())
}

func webrtcmock() {
	opts := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", true),
		chromedp.Flag("ignore-certificate-errors", true),
		chromedp.Flag("unsafely-treat-insecure-origin-as-secure", "http://192.168.3.250"),
	}

	ctx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx2, cancel2 := chromedp.NewContext(ctx)
	defer cancel2()

	err := chromedp.Run(ctx2, chromedp.Navigate("http://192.168.3.250"))

	// run task list
	err = chromedp.Run(ctx2,
		chromedp.SendKeys("#app > div > section > main > div > div:nth-child(2) > div > div:nth-child(2) > input", "642e1fe27d6011dd5db7d774"),
	)
	if err != nil {
		log.Fatal(err)
	}

	err = chromedp.Run(ctx2,
		chromedp.Click("#app > div > section > main > div > div:nth-child(2) > div > button.el-button.el-button--primary > span"),
	)

	err = chromedp.Run(ctx2, chromedp.Sleep(5*time.Hour))
	if err != nil {
		panic(err)
	}
}
