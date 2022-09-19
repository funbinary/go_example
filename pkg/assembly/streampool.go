package assembly

import (
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/google/gopacket/layers"
	"sync"
	"time"
)

const initialAllocSize = 1024

type StreamPool struct {
	conns              map[key]*connection
	users              int //
	mu                 sync.RWMutex
	factory            StreamFactory
	free               []*connection
	all                [][]connection
	nextAlloc          int
	newConnectionCount int64
}

func NewStreamPool(factory StreamFactory) *StreamPool {
	return &StreamPool{
		conns:     make(map[key]*connection, initialAllocSize),
		free:      make([]*connection, 0, initialAllocSize),
		factory:   factory,
		nextAlloc: initialAllocSize,
	}
}

func (self *StreamPool) GetConnection(k key, end bool, ts time.Time, tcp *layers.TCP, ac AssemblerContext) (*connection, *halfconnection, *halfconnection) {
	self.mu.RLock()
	conn, half, rev := self.getHalf(k)
	self.mu.RUnlock()
	// 找到直接返回
	if end || conn != nil {
		return conn, half, rev
	}
	// 创建stream
	s := self.factory.New(k[0], k[1], tcp, ac)
	self.mu.Lock()
	defer self.mu.Unlock()
	// 创建连接
	conn, half, rev = self.newConnection(k, s, ts)
	self.conns[k] = conn
	return conn, half, rev
}

func (p *StreamPool) getHalf(k key) (*connection, *halfconnection, *halfconnection) {
	conn := p.conns[k]
	if conn != nil {
		return conn, &conn.c2s, &conn.s2c
	}
	rk := k.Reverse()
	conn = p.conns[rk]
	if conn != nil {
		return conn, &conn.s2c, &conn.c2s
	}
	return nil, nil, nil
}

func (p *StreamPool) newConnection(k key, s Stream, ts time.Time) (c *connection, h *halfconnection, r *halfconnection) {
	if memlog {
		p.newConnectionCount++
		if p.newConnectionCount&0x7FFF == 0 {
			log.Debugf("StreamPool:%d requests,%d used,%d free", p.newConnectionCount, len(p.conns), len(p.free))
		}
	}
	if len(p.free) == 0 {
		p.grow()
	}
	index := len(p.free) - 1
	c, p.free = p.free[index], p.free[:index]
	c.reset(k, s, ts)
	return c, &c.c2s, &c.s2c
}

func (p *StreamPool) grow() {
	conns := make([]connection, p.nextAlloc)
	p.all = append(p.all, conns)
	for i := range conns {
		p.free = append(p.free, &conns[i])
	}
	if memlog {
		log.Debugf("StreamPool: created:%d new connections", p.nextAlloc)
	}
	p.nextAlloc *= 2
}

// Dump logs all connections
func (p *StreamPool) Dump() {
	p.mu.Lock()
	defer p.mu.Unlock()
	log.Debugf("Remaining %d connections: ", len(p.conns))
	for _, conn := range p.conns {
		log.Debugf("%v %s", conn.key, conn)
	}
}

func (p *StreamPool) Remove(conn *connection) {
	p.mu.Lock()
	if _, ok := p.conns[conn.key]; ok {
		delete(p.conns, conn.key)
		p.free = append(p.free, conn)
	}
	p.mu.Unlock()
}
