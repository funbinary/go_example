package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/funbinary/go_example/example/crawler/08/collect"
)

// xpath
func main() {
	url := "https://book.douban.com/subject/1007305/"
	fetch := collect.BrowserFetch{}
	body, err := fetch.Get(url)

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("htmlquery.Parse failed:%v", err)
	}
	doc.Find("div.small_cardcontent__BTALp h2").Each(func(i int, selection *goquery.Selection) {
		title := selection.Text()
		fmt.Println("fetch card:", i, " ", title)
	})

}
