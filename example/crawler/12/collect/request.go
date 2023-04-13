package collect

type Request struct {
	Url       string // 访问的防战
	Cookie    string
	ParseFunc func([]byte, *Request) ParseResult // 解析从网站获取到的网站信息
}

type ParseResult struct {
	Requesrts []*Request
	Items     []interface{}
}
