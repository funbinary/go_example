package nfsv2

import (
	"encoding/hex"
	"fmt"
	"github.com/funbinary/go_example/pkg/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

type GetAttrCall struct {
	FileHandle []byte `xdr:"limit=32"`
}

func (self *GetAttrCall) String() string {
	var s strings.Builder
	s.WriteString(hex.EncodeToString(self.FileHandle))
	return s.String()
}

type GetAttrReply struct {
	Nfs2Status
	Fattr
}

func (self *GetAttrReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("attributes:\n", self.Fattr.String()))
	return s.String()
}

type SetAttrCall struct {
	FileHandle []byte `xdr:"limit=32"`
	Attr       Sattr
}

func (self *SetAttrCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("File Handle:", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("new attributes:\n", self.Attr.String()))
	return s.String()
}

type SetAttrReply struct {
	Nfs2Status
	Fattr
}

func (self *SetAttrReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("attributes:\n", self.Fattr.String()))
	return s.String()
}

type LookUpCall struct {
	DirOpArgs
}

func (self *LookUpCall) String() string {
	var s strings.Builder
	s.WriteString(self.DirOpArgs.String())
	return s.String()
}

type LookUpReply struct {
	Nfs2Status
	DirOpRes
}

func (self *LookUpReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("File Handle:", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("DirAttr:\n", self.Attr.String()))
	return s.String()
}

type ReadLinkCall struct {
	FileHandle []byte `xdr:"limit=32"`
}

func (self *ReadLinkCall) String() string {
	var s strings.Builder
	s.WriteString(hex.EncodeToString(self.FileHandle))
	return s.String()
}

type ReadLinkReply struct {
	Nfs2Status
	FilePath string
}

func (self *ReadLinkReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("File Path:", self.FilePath))
	return s.String()
}

type ReadCall struct {
	FileHandle []byte `xdr:"limit=32"`
	Offset     uint32
	Count      uint32
	TotalCount uint32 //没有使用
}

func (self *ReadCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("Offset:", self.Offset))
	s.WriteString(fmt.Sprintln("Count:", self.Count))
	s.WriteString(fmt.Sprintln("Total Count:", self.TotalCount))
	return s.String()
}

type ReadReply struct {
	Nfs2Status
	Attr Fattr
	Data []byte
}

func (self *ReadReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("obj attributes:\n", self.Attr.String()))
	s.WriteString(fmt.Sprintln("data lenth:", len(self.Data)))
	return s.String()
}

type WriteCall struct {
	FileHandle  []byte `xdr:"limit=32"`
	BeginOffset uint32
	Offset      uint32
	TotalCount  uint32
	Data        []byte
}

func (self *WriteCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("BeginOffset:", self.BeginOffset))
	s.WriteString(fmt.Sprintln("Offset:", self.Offset))
	s.WriteString(fmt.Sprintln("Total Count:", self.TotalCount))
	s.WriteString(fmt.Sprintln("Real Count:", len(self.Data)))
	return s.String()
}

type WriteReply struct {
	Nfs2Status
	Attr Fattr
}

func (self *WriteReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("obj attributes:\n", self.Attr.String()))
	return s.String()
}

type CreateCall struct {
	DirOpArgs
	Attr Sattr
}

func (self *CreateCall) String() string {
	var s strings.Builder
	s.WriteString(self.DirOpArgs.String())
	s.WriteString(fmt.Sprintln("Set attributes:\n", self.Attr.String()))

	return s.String()
}

type CreateReply struct {
	Nfs2Status
	DirOpRes
}

func (self *CreateReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("File Handle:", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("DirAttr:\n", self.Attr.String()))
	return s.String()
}

type RemoveCall struct {
	DirOpArgs
}

func (self *RemoveCall) String() string {
	var s strings.Builder
	s.WriteString(self.DirOpArgs.String())
	return s.String()
}

type RemoveReply struct {
	Nfs2Status
}

func (self *RemoveReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	return s.String()
}

type RenameCall struct {
	From DirOpArgs
	To   DirOpArgs
}

func (self *RenameCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("from:\n ", self.From))
	s.WriteString(fmt.Sprintln("to:\n ", self.To))
	return s.String()
}

type RenameReply struct {
	Nfs2Status
}

func (self *RenameReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	return s.String()
}

type LinkCall struct {
	FileHandle []byte `xdr:"limit=32"` //from
	To         DirOpArgs
}

func (self *LinkCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("from:", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("to:\n ", self.To))
	return s.String()
}

type LinkReply struct {
	Nfs2Status
}

func (self *LinkReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	return s.String()
}

type SymlinkCall struct {
	From DirOpArgs
	To   string
	Attr Sattr
}

func (self *SymlinkCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("from:", self.From))
	s.WriteString(fmt.Sprintln("to: ", self.To))
	s.WriteString(fmt.Sprintln("set attr:\n ", self.Attr.String()))
	return s.String()
}

type SymlinkReply struct {
	Nfs2Status
}

func (self *SymlinkReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	return s.String()
}

type MkdirCall struct {
	DirOpArgs
	Attr Sattr
}

func (self *MkdirCall) String() string {
	var s strings.Builder
	s.WriteString(self.DirOpArgs.String())
	s.WriteString(fmt.Sprintln("Set attributes:\n", self.Attr.String()))

	return s.String()
}

type MkdirReply struct {
	Nfs2Status
	DirOpRes
}

func (self *MkdirReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("File Handle:", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("DirAttr:\n", self.Attr.String()))
	return s.String()
}

type RmdirCall struct {
	DirOpArgs
}

func (self *RmdirCall) String() string {
	var s strings.Builder
	s.WriteString(self.DirOpArgs.String())
	return s.String()
}

type RmdirReply struct {
	Nfs2Status
}

func (self *RmdirReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	return s.String()
}

type ReadDirCall struct {
	FileHandle []byte `xdr:"limit=32"`
	Cookie     uint32
	Count      uint32
}

func (self *ReadDirCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("Cookie:", self.Cookie))
	s.WriteString(fmt.Sprintln("Count:", self.Count))
	return s.String()
}

type Entry struct {
	Field    uint32
	Filename []byte
	Cookie   uint32
}

type Entry2 struct {
	IsSet bool  `xdr:"union"`
	Entry Entry `xdr:"unioncase=1"`
}

type ReadDirReply struct {
	Nfs2Status
	Entries []*Entry
}

func (self *ReadDirReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("entrys:"))
	for _, v := range self.Entries {
		s.WriteString(fmt.Sprintln("    field:", v.Field))
		switch charset.GetStrCoding(v.Filename) {
		case charset.GBK:
			utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(v.Filename) //将gbk再转换为utf-8
			s.WriteString(fmt.Sprintln("    file name:", string(utf8Data)))
		default:
			s.WriteString(fmt.Sprintln("    file name:", string(v.Filename)))
		}
		s.WriteString(fmt.Sprintln("    cookie:", v.Cookie))
	}

	return s.String()
}

type StatFsCall struct {
	FileHandle []byte `xdr:"limit=32"`
}

func (self *StatFsCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	return s.String()
}

type StatFsReply struct {
	Nfs2Status
	TSize  uint32 //用字节表示的最优化的传输尺寸。这是服务器在READ 和 WRITE请求中的 最想要的数据字节数。
	BSize  uint32 //文件系统用字节表示的块尺寸。.
	Blocks uint32 //文件系统中 "bsize"块的总数。
	BFree  uint32 //文件系统中自由的“bsize”块的数目。
	BAvail uint32 //无特权用户可用的"bsize"块的数目。
}

func (self *StatFsReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.Nfs2Status))
	s.WriteString(fmt.Sprintln("TSize:", self.TSize))
	s.WriteString(fmt.Sprintln("BSize:", self.BSize))
	s.WriteString(fmt.Sprintln("Blocks:", self.Blocks))
	s.WriteString(fmt.Sprintln("BFree:", self.BFree))
	s.WriteString(fmt.Sprintln("BAvail:", self.BAvail))
	return s.String()
}
