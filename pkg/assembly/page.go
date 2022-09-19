package assembly

import (
	"fmt"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/google/gopacket"
	"sync"
	"time"
)

const pageBytes = 1900

/* page: implements a byteContainer */

// page is used to store TCP data we're not ready for yet (out-of-order
// packets).  Unused pages are stored in and returned from a pageCache, which
// avoids memory allocation.  Used pages are stored in a doubly-linked list in
// a connection.
type page struct {
	bytes      []byte
	seq        Sequence
	prev, next *page
	buf        [pageBytes]byte
	ac         AssemblerContext // only set for the first page of a packet
	seen       time.Time
	start, end bool
}

func (p *page) getBytes() []byte {
	return p.bytes
}
func (p *page) captureInfo() gopacket.CaptureInfo {
	return p.ac.GetCaptureInfo()
}
func (p *page) assemblerContext() AssemblerContext {
	return p.ac
}
func (p *page) convertToPages(pc *pageCache, skip int, ac AssemblerContext) (*page, *page, int) {
	if skip != 0 {
		p.bytes = p.bytes[skip:]
		p.seq = p.seq.Add(skip)
	}
	p.prev, p.next = nil, nil
	return p, p, 1
}
func (p *page) length() int {
	return len(p.bytes)
}
func (p *page) release(pc *pageCache) int {
	pc.replace(p)
	return 1
}
func (p *page) isStart() bool {
	return p.start
}
func (p *page) isEnd() bool {
	return p.end
}
func (p *page) getSeq() Sequence {
	return p.seq
}
func (p *page) isPacket() bool {
	return p.ac != nil
}
func (p *page) String() string {
	return fmt.Sprintf("page@%p{seq: %v, bytes:%d, -> nextSeq:%v} (prev:%p, next:%p)", p, p.seq, len(p.bytes), p.seq+Sequence(len(p.bytes)), p.prev, p.next)
}

/*
 * pageCache
 */
// pageCache is a concurrency-unsafe store of page objects we use to avoid
// memory allocation as much as we can.
type pageCache struct {
	pagePool     *sync.Pool
	used         int
	pageRequests int64
}

func newPageCache() *pageCache {
	pc := &pageCache{
		pagePool: &sync.Pool{
			New: func() interface{} { return new(page) },
		}}
	return pc
}

// next returns a clean, ready-to-use page object.
func (c *pageCache) next(ts time.Time) (p *page) {
	if memlog {
		c.pageRequests++
		if c.pageRequests&0xFFFF == 0 {
			log.Debugf("PageCache:%v requested:%v used,", c.pageRequests, c.used)
		}
	}

	p = c.pagePool.Get().(*page)
	p.seen = ts
	p.bytes = p.buf[:0]
	c.used++
	if memlog {
		log.Debugf("allocator returns %s\n", p)
	}

	return p
}

// replace replaces a page into the pageCache.
func (c *pageCache) replace(p *page) {
	c.used--
	if memlog {
		log.Debugf("replacing %s\n", p)
	}
	p.prev = nil
	p.next = nil
	c.pagePool.Put(p)
}
