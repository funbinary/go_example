package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/funbinary/go_example/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"net/http"
)

// xpath
func main() {
	url := "https://www.thepaper.cn/"
	body, err := Fetch(url)

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}
	doc, err := htmlquery.Parse(bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("htmlquery.Parse failed:%v", err)
	}
	nodes := htmlquery.Find(doc, `//div[@class="small_cardcontent__BTALp"]//h2`)
	for _, node := range nodes {
		fmt.Println("fetch card news:", node.FirstChild.Data)
	}
}

func Fetch(url string) ([]byte, error) {
	var e encoding.Encoding
	resp, err := http.Get(url)

	if err != nil {
		return nil, errors.Wrap(err, "fetch url error")
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("Error status code:%v", resp.StatusCode)
	}

	r := bufio.NewReader(resp.Body)
	e, err = DeterminEncoding(r)
	utf8r := transform.NewReader(r, e.NewDecoder())
	return io.ReadAll(utf8r)
}

func DeterminEncoding(r *bufio.Reader) (encoding.Encoding, error) {
	bytes, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8, err
	}

	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e, nil
}
