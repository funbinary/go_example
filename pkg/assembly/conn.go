package assembly

import (
	"fmt"
	"sync"
	"time"
)

/* one-way connection, i.e. halfconnection */
type halfconnection struct {
	dir         TCPFlowDirection
	pages       int      // Number of pages used (both in first/last and saved)
	saved       *page    // Doubly-linked list of in-order pages (seq < nextSeq) already given to Stream who told us to keep
	first, last *page    // Doubly-linked list of out-of-order pages (seq > nextSeq)
	nextSeq     Sequence // sequence number of in-order received bytes
	ackSeq      Sequence
	created     time.Time
	lastTime    time.Time //最后一个包的接受时间，初始值为当前包的接受时间
	stream      Stream
	closed      bool
	// for stats
	queuedBytes    int
	queuedPackets  int
	overlapBytes   int
	overlapPackets int
}

func (self *halfconnection) SetLastPacketRecvTime(t time.Time) {
	if self.lastTime.Before(t) {
		self.lastTime = t
	}
}

func (half *halfconnection) String() string {
	closed := ""
	if half.closed {
		closed = "closed "
	}
	return fmt.Sprintf("%screated:%v, last:%v", closed, half.created, half.lastTime)
}

// Dump returns a string (crypticly) describing the halfconnction
func (half *halfconnection) Dump() string {
	s := fmt.Sprintf("pages: %d\n"+
		"nextSeq: %d\n"+
		"ackSeq: %d\n"+
		"Seen :  %s\n"+
		"dir:    %s\n", half.pages, half.nextSeq, half.ackSeq, half.lastTime, half.dir)
	nb := 0
	for p := half.first; p != nil; p = p.next {
		s += fmt.Sprintf("	Page[%d] %s len: %d\n", nb, p, len(p.bytes))
		nb++
	}
	return s
}

// 双向连接
type connection struct {
	key      key // client->server
	c2s, s2c halfconnection
	mu       sync.Mutex
}

func (c *connection) reset(k key, s Stream, ts time.Time) {
	c.key = k
	base := halfconnection{
		nextSeq:  invalidSequence,
		ackSeq:   invalidSequence,
		created:  ts,
		lastSeen: ts,
		stream:   s,
	}
	c.c2s, c.s2c = base, base
	c.c2s.dir, c.s2c.dir = TCPDirClientToServer, TCPDirServerToClient
}

func (c *connection) lastSeen() time.Time {
	if c.c2s.lastSeen.Before(c.s2c.lastSeen) {
		return c.s2c.lastSeen
	}

	return c.c2s.lastSeen
}

func (c *connection) String() string {
	return fmt.Sprintf("c2s: %s, s2c: %s", &c.c2s, &c.s2c)
}
