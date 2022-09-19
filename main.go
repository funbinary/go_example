package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/httputil"
	"net/url"
)

/*
   gin框架实现文件下载功能
*/

//主函数
func main() {
	e := gin.Default()
	v1 := e.Group("/v1")
	v1.POST("/startRecord", ProxyRecord)
	v1.POST("/stopRecord", ProxyRecord)
	v1.POST("/concat", ProxyRecord)
	e.Run(":9001")

}

type DebugTransport struct{}

func (DebugTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	b, err := httputil.DumpRequestOut(r, false)
	if err != nil {
		return nil, err
	}
	logrus.Debugf(string(b))
	return http.DefaultTransport.RoundTrip(r)
}

func ProxyRecord(c *gin.Context) {
	var proxyUrl = new(url.URL)
	proxyUrl.Scheme = "http"
	proxyUrl.Host = "192.168.3.250" + ":" + "9004"
	logrus.Debugf(proxyUrl.Host)
	proxy := httputil.NewSingleHostReverseProxy(proxyUrl)
	proxy.Transport = DebugTransport{}
	proxy.ServeHTTP(c.Writer, c.Request)
	c.Abort()
}
