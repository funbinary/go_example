package collect

import (
	"github.com/funbinary/go_example/pkg/errors"
	"time"
)

type Request struct {
	Url       string // 访问的防战
	Cookie    string
	WaitTime  time.Duration
	Depth     int
	MaxDepth  int
	ParseFunc func([]byte, *Request) ParseResult // 解析从网站获取到的网站信息
}

func (r *Request) Check() error {
	if r.Depth < r.MaxDepth {
		return errors.New("Max depth limit reached")
	}
	return nil
}

type ParseResult struct {
	Requesrts []*Request
	Items     []interface{}
}
