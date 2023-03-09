package main

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

type HTTPProxy struct {
	proxy *httputil.ReverseProxy
}

func NewHTTPProxy(target string) (*HTTPProxy, error) {
	u, err := url.Parse(target)
	if err != nil {
		return nil, err
	}
	return &HTTPProxy{httputil.NewSingleHostReverseProxy(u)}, nil
}

func (h *HTTPProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.proxy.ServeHTTP(w, r)
}

func main() {
	proxy, err := NewHTTPProxy("http://192.168.3.250:80")
	if err != nil {

	}
	http.Handle("/", proxy)
	http.ListenAndServe(":8081", nil)
}
