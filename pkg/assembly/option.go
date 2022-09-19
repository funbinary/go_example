package assembly

type option func(*Assembler)

// WithMaxBufferedPagesTotal
// 等待无序包时要缓冲的page总数最大值
// 一旦达到这个上限值， Assembler将会降级刷新每个连接的
// 如果<=0将被忽略。
func WithMaxBufferedPagesTotal(m int) option {
	return func(a *Assembler) {
		a.maxBufferedPagesTotal = m
	}
}

// WithMaxBufferedPagesPerConnection
// 单个连接缓冲的page最大值
// 如果达到上限，则将刷新最小序列号以及任何连续数据。
// 如果<= 0，这将被忽略。
func WithMaxBufferedPagesPerConnection(m int) option {
	return func(a *Assembler) {
		a.maxBufferedPagesPerConnection = m
	}
}
