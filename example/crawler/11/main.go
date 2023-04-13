package main

import (
	"github.com/funbinary/go_example/example/crawler/10/collect"
	"github.com/funbinary/go_example/example/crawler/10/proxy"
	"github.com/funbinary/go_example/example/crawler/11/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"time"
)

// xpath
func main() {
	plugin, c := log.NewFilePlugin("./log.txt", zapcore.InfoLevel)
	defer c.Close()
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	proxyURLs := []string{"http://127.0.0.1:10809", "http://127.0.0.1:10809"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}
	url := "https://google.com"
	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}

	body, err := f.Get(url)
	if err != nil {
		logger.Error("read content failed",
			zap.Error(err),
		)
		return
	}
	logger.Info("get content", zap.Int("len", len(body)))

}
