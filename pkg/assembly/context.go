package assembly

import "github.com/google/gopacket"

// AssemblerContext provides method to get metadata
type AssemblerContext interface {
	GetCaptureInfo() gopacket.CaptureInfo
}

// Implements AssemblerContext for Assemble()
type assemblerSimpleContext gopacket.CaptureInfo

func (asc *assemblerSimpleContext) GetCaptureInfo() gopacket.CaptureInfo {
	return gopacket.CaptureInfo(*asc)
}
