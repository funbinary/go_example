package main

import (
	"bufio"
	"fmt"
	"github.com/funbinary/go_example/pkg/errors"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
	"io"
	"net/http"
	"regexp"
)

// var headerRe = regexp.MustCompile(`<div class="news_li"[\s\S]*?<h2>[\s\S]*?<a.*?target="_blank">([\s\S]*?)</a>`)
var headerRe = regexp.MustCompile(`<div class="small_cardcontent__BTALp"[\s\S]*?<h2>([\s\S]*?)</h2>`)

func main() {
	url := "https://www.thepaper.cn/"
	body, err := Fetch(url)

	if err != nil {
		fmt.Println("read content failed:%v", err)
		return
	}
	//fmt.Println(string(body))
	matches := headerRe.FindAllSubmatch(body, -1)
	for _, m := range matches {
		fmt.Println("fetch card news:", string(m[1]))
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
