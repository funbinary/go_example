package main

import (
	"fmt"
	"github.com/funbinary/go_example/example/crawler/12/collect"
	"github.com/funbinary/go_example/example/crawler/12/log"
	"github.com/funbinary/go_example/example/crawler/12/parse/doubangroup"
	"github.com/funbinary/go_example/example/crawler/12/proxy"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"time"
)

// xpath
func main() {
	plugin := log.NewStdoutPlugin(zapcore.InfoLevel)
	logger := log.NewLogger(plugin)
	logger.Info("log init end")

	proxyURLs := []string{"http://127.0.0.1:10809", "http://127.0.0.1:10809"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		logger.Error("RoundRobinProxySwitcher failed")
	}

	cookie := `ll="118201"; __utmc=30149280; ck=QaGr; push_noty_num=0; push_doumail_num=0; __utmv=30149280.21545; __yadk_uid=CY4XlZtUkKWowjb53K8SISQTgqj8YOOU; douban-fav-remind=1; frodotk_db="8df2541269e216dca9d6fc373da64494"; bid=dPuzdR0mG9M; gr_user_id=690ec6c6-4e7f-4277-b959-b829fd4aef5a; viewed="1007305_1475839_25913349"; __gads=ID=613f831a31c6ac24-225718cbcadc0032:T=1679924466:RT=1679924466:S=ALNI_MaDEdHHhIEtazV6BqOobp1mDpI4Ug; __gpi=UID=00000be220b6ea7c:T=1679924466:RT=1679924466:S=ALNI_Mbp-472jjdHsL0xjpHPnuuWAacAEg; _pk_ref.100001.8cb4=["","",1680704148,"https://time.geekbang.org/column/article/612328"]; _pk_ses.100001.8cb4=*; ap_v=0,6.0; __utma=30149280.1773533084.1677888507.1679440673.1680704158.4; __utmz=30149280.1680704158.4.2.utmcsr=time.geekbang.org|utmccn=(referral)|utmcmd=referral|utmcct=/column/article/612328; __utmt=1; _pk_id.100001.8cb4=bb24eb830bd259ee.1677888506.8.1680704178.1679440672.; __utmb=30149280.12.5.1680704178246`

	var worklist []*collect.Request
	for i := 0; i <= 0; i += 25 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		worklist = append(worklist, &collect.Request{
			Url:       str,
			Cookie:    cookie,
			ParseFunc: doubangroup.ParseURL,
		})
	}

	var f collect.Fetcher = &collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}

	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			body, err := f.Get(item)
			time.Sleep(1 * time.Second)
			if err != nil {
				logger.Error("read content failed",
					zap.Error(err),
				)
				continue
			}
			res := item.ParseFunc(body, item)
			for _, item := range res.Items {
				logger.Info("result",
					zap.String("get url:", item.(string)))
			}
			worklist = append(worklist, res.Requesrts...)
		}
	}

}
