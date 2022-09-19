package assembly

import (
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

var (
	memlog = false
)

const assemblerReturnValueInitialSize = 16

// byteContainer is either a page or a livePacket
type byteContainer interface {
	getBytes() []byte
	length() int
	convertToPages(*pageCache, int, AssemblerContext) (*page, *page, int)
	captureInfo() gopacket.CaptureInfo
	assemblerContext() AssemblerContext
	release(*pageCache) int
	isStart() bool
	isEnd() bool
	getSeq() Sequence
	isPacket() bool
}

type Assembler struct {
	maxBufferedPagesTotal         int // 等待无序包时要缓冲的page总数最大值
	maxBufferedPagesPerConnection int // 单个连接缓冲的page最大值
	ret                           []byteContainer
	pc                            *pageCache
	connPool                      *StreamPool
	//cacheLP                       livePacket
	//cacheSG                       reassemblyObject
	start bool
}

func NewAssembler(pool *StreamPool, opts ...option) *Assembler {
	a := &Assembler{
		ret:      make([]byteContainer, 0, assemblerReturnValueInitialSize),
		pc:       newPageCache(),
		connPool: pool,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (self *Assembler) AssembleWithContext(netFlow gopacket.Flow, t *layers.TCP, ac AssemblerContext) {
	// 4元组key
	key := key{netFlow, t.TransportFlow()}
	cinfo := ac.GetCaptureInfo()
	// 连接/当前的连接/相反的连接
	conn, half, rev := self.connPool.GetConnection(key, false, cinfo.Timestamp, t, ac)
	if conn == nil {
		log.Debugf("%v got empty packet on otherwise empty connection", key)
		return
	}
	// 设置最后一个包的接收时间
	half.SetLastPacketRecvTime(cinfo.Timestamp)

}
