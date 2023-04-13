package engine

import (
	"github.com/funbinary/go_example/example/crawler/13/collect"
	"go.uber.org/zap"
)

type ScheduleEngine struct {
	requestCh chan *collect.Request //负责接收请求
	workerCh  chan *collect.Request //负责分配任务给 worker
	WorkCount int                   //为执行任务的数量，可以灵活地去配置。
	Fetcher   collect.Fetcher
	Logger    *zap.Logger
	out       chan collect.ParseResult ////负责处理爬取后的数据，完成下一步的存储操作。schedule 函数会创建调度程序，负责的是调度的核心逻辑。
	Seeds     []*collect.Request
}

func (s *ScheduleEngine) Run() {
	s.requestCh = make(chan *collect.Request)
	s.workerCh = make(chan *collect.Request)
	s.out = make(chan collect.ParseResult)
	go s.Schedule()
	// 创建指定数量的 worker，完成实际任务的处理
	// 其中
	for i := 0; i < s.WorkCount; i++ {
		go s.CreateWork()
	}
	s.HandleResult()
}

func (s *ScheduleEngine) Schedule() {
	// workerCh
	var reqQueue = s.Seeds
	go func() {
		for {
			var req *collect.Request
			var ch chan *collect.Request

			//如果任务队列 reqQueue 大于 0，意味着有爬虫任务，这时我们获取队列中第一个任务，并将其剔除出队列。

			if len(reqQueue) > 0 {
				req = reqQueue[0]
				reqQueue = reqQueue[1:]
				ch = s.workerCh
			}
			select {
			case r := <-s.requestCh:
				// 接收来自外界的请求，并将请求存储到 reqQueue 队列中
				reqQueue = append(reqQueue, r)

			case ch <- req:
				// ch <- req 会将任务发送到 workerCh 通道中，等待 worker 接收。
			}
		}
	}()

}

func (s *ScheduleEngine) CreateWork() {
	for {
		// 接收到调度器分配的任务；
		r := <-s.workerCh
		// 访问服务器
		body, err := s.Fetcher.Get(r)
		if err != nil {
			s.Logger.Error("can't fetch ",
				zap.Error(err),
			)
			continue
		}
		//解析服务器返回的数据
		result := r.ParseFunc(body, r)
		// 将返回的数据发送到 out 通道中，方便后续的处理。
		s.out <- result
	}
}

func (s *ScheduleEngine) HandleResult() {
	for {
		select {
		// 接收所有 worker 解析后的数据
		case result := <-s.out:
			// 要进一步爬取的 Requests 列表将全部发送回 s.requestCh 通道
			for _, req := range result.Requesrts {
				s.requestCh <- req
			}
			//包含了我们实际希望得到的结果，所以我们先用日志把结果打印出来
			for _, item := range result.Items {
				// todo: store
				s.Logger.Sugar().Info("get result", item)
			}
		}
	}
}
