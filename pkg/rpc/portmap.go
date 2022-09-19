package rpc

type PortmapProcedure uint32

// PortmapProcedure
const (
	PortmapProcNull    PortmapProcedure = 0
	PortmapProcGetPort PortmapProcedure = 3
)

func (m PortmapProcedure) String() string {
	switch m {
	case PortmapProcNull:
		return "Null"
	case PortmapProcGetPort:
		return "GetPort"
	default:
		return "Unknown"
	}
}
