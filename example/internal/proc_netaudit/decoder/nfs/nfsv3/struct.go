package nfsv3

import (
	"encoding/hex"
	"fmt"
	"github.com/funbinary/go_example/pkg/charset"
	"golang.org/x/text/encoding/simplifiedchinese"
	"strings"
)

type GetAttrCall struct {
	FileHandle []byte
}

func (self *GetAttrCall) String() string {
	var s strings.Builder
	s.WriteString(hex.EncodeToString(self.FileHandle))

	return s.String()
}

type GetAttrReply struct {
	NFSStatus
	Fattr
}

func (self *GetAttrReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	s.WriteString(fmt.Sprintln("attributes:\n", self.Fattr.String()))
	return s.String()
}

type SetAttrCall struct {
	FileHandle []byte
	Attr       Sattr3
	Guard      Sattrguard3
}

func (self *SetAttrCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("FileHandle:", hex.EncodeToString(self.FileHandle)))
	if self.Attr.HasSet() {
		s.WriteString(fmt.Sprintln("new attributes:\n", self.Attr.String()))
	}
	//if self.Guard.HasCheck() {
	//	s.WriteString(fmt.Sprintln("guard:\n", self.Guard.String()))
	//}

	return s.String()
}

type SetAttrReply struct {
	NFSStatus
	Wcc
}

func (self *SetAttrReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.Wcc.HasSet() {
		s.WriteString(fmt.Sprintln("wcc:\n", self.Wcc))
	}

	return s.String()
}

type LookUpCall struct {
	DirOpArg
}

func (self *LookUpCall) String() string {
	var s strings.Builder
	s.WriteString(self.DirOpArg.String())

	return s.String()
}

type LookUpReply struct {
	NFSStatus
	FH      []byte
	Attr    PostOpAttr
	DirAttr PostOpAttr
}

func (self *LookUpReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	switch self.NFSStatus {
	case NFSStatusOk:
		s.WriteString(fmt.Sprintln("FileHandle:", hex.EncodeToString(self.FH)))
		if self.Attr.HasSet() {
			s.WriteString(fmt.Sprintln("Attr:\n", self.Attr.String()))
		}
		if self.DirAttr.HasSet() {
			s.WriteString(fmt.Sprintln("DirAttr:\n", self.DirAttr.String()))
		}
	default:
		if self.Attr.HasSet() {
			s.WriteString(fmt.Sprintln("Attr:\n", self.Attr.String()))
		}
	}
	return s.String()
}

type AccessCall struct {
	FileHandle []byte // root handle
	Access            //todo
}

func (self *AccessCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("access:", self.Access))

	return s.String()
}

type AccessReply struct {
	NFSStatus
	Attr   PostOpAttr
	Access //todo
}

func (self *AccessReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.Attr.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.Attr.String()))
	}

	s.WriteString(fmt.Sprintln("access:", self.Access))

	return s.String()
}

type ReadLinkCall struct {
	FileHandle []byte
}

func (self *ReadLinkCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))

	return s.String()
}

type ReadLinkReply struct {
	NFSStatus
	Attr PostOpAttr
	Path string
}

func (self *ReadLinkReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.Attr.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.Attr.String()))
	}
	s.WriteString(fmt.Sprintf("path:", self.Path))

	return s.String()
}

type ReadCall struct {
	FileHandle []byte
	Offset     uint64
	Count      uint32
}

func (self *ReadCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("Offset:", self.Offset))
	s.WriteString(fmt.Sprintln("Count:", self.Count))

	return s.String()
}

type ReadReply struct {
	NFSStatus
	Attr  PostOpAttr
	Count uint32
	EOF   uint32
	Data  []byte
}

func (self *ReadReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.Attr.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.Attr.String()))
	}

	s.WriteString(fmt.Sprintln("count:", self.Count))
	s.WriteString(fmt.Sprintln("data lenth:", len(self.Data)))

	return s.String()
}

type WriteCall struct {
	Handle []byte
	Offset uint64
	Count  uint32
	How    uint32 //todo :format 2-filesync
	Data   []byte
}

func (self *WriteCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.Handle)))
	s.WriteString(fmt.Sprintln("offset: ", self.Offset))
	s.WriteString(fmt.Sprintln("count: ", self.Count))
	s.WriteString(fmt.Sprintln("how: ", self.How))
	s.WriteString(fmt.Sprintln("data length: ", len(self.Data)))
	if len(self.Data) > 0 {
		s.WriteString(fmt.Sprintln(hex.EncodeToString(self.Data)))
	}
	return s.String()
}

type WriteReply struct {
	NFSStatus
	Wcc
	Count    uint32
	Commit   uint32
	Verifier uint64
}

func (self *WriteReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	s.WriteString(fmt.Sprintln("wcc:\n", self.Wcc.String()))
	s.WriteString(fmt.Sprintln("count:", self.Count))
	s.WriteString(fmt.Sprintln("commit:", self.Commit))
	s.WriteString(fmt.Sprintln("verifier:", self.Verifier))

	return s.String()
}

type CreateCall struct {
	DirOpArg
	// 0 : UNCHECKED (default)
	// 1 : GUARDED
	// 2 : EXCLUSIVE
	Mode uint32
	Attr Sattr3
}

func (self *CreateCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("where:"))
	s.WriteString(self.DirOpArg.String())
	s.WriteString("create mode:")
	switch self.Mode {
	case 0:
		s.WriteString("UNCHECKED")
	case 1:
		s.WriteString("GUARDED")
	case 2:
		s.WriteString("EXCLUSIVE")
	}
	s.WriteString("\n")
	if self.Attr.HasSet() {
		s.WriteString(self.Attr.String())
	}
	return s.String()
}

type CreateReply struct {
	NFSStatus
	FileHandle PostOpFH3
	Attr       PostOpAttr
	Wcc
}

func (self *CreateReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.FileHandle.IsSet {
		s.WriteString(fmt.Sprintln("File handle:", self.FileHandle))
	}

	if self.Attr.HasSet() {
		s.WriteString(fmt.Sprintln("attributes:"))
		s.WriteString(self.Attr.String())
	}
	s.WriteString(fmt.Sprintln("wcc:\n", self.Wcc.String()))

	return s.String()
}

type MkdirCall struct {
	Where DirOpArg
	Attrs Sattr3
}

func (self *MkdirCall) String() string {
	var s strings.Builder
	s.WriteString(self.Where.String())
	if self.Attrs.HasSet() {
		s.WriteString(self.Attrs.String())
	}
	return s.String()
}

type MkdirReply struct {
	NFSStatus
	FH   PostOpFH3
	Attr PostOpAttr
	Wcc
}

func (self *MkdirReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.FH.IsSet {
		s.WriteString(fmt.Sprintln("File handle:", self.FH))
	}

	if self.Attr.HasSet() {
		s.WriteString(fmt.Sprintln("attributes:"))
		s.WriteString(self.Attr.String())
	}
	s.WriteString(fmt.Sprintln("wcc:\n", self.Wcc.String()))

	return s.String()
}

type SymlinkCall struct {
	DirOpArg
	SymLinkData3
}

func (self *SymlinkCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("Dir op arg:\n", self.DirOpArg))
	s.WriteString(fmt.Sprintln("Symlink data:\n ", self.SymLinkData3))

	return s.String()
}

type SymlinkReply struct {
	NFSStatus
	Obj     PostOpFH3
	ObjAttr PostOpAttr
	Wcc
}

func (self *SymlinkReply) String() string {

	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.Obj.IsSet {
		s.WriteString(fmt.Sprintln("File handle:", self.Obj))
	}

	if self.ObjAttr.HasSet() {
		s.WriteString(fmt.Sprintln("attributes:"))
		s.WriteString(self.ObjAttr.String())
	}
	s.WriteString(fmt.Sprintf("wcc:\n", self.Wcc))

	return s.String()
}

type MknodCall struct {
	DirOpArg
	FileType
	Sattr3
}

func (self *MknodCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("Dir op arg:\n ", self.DirOpArg))
	s.WriteString(fmt.Sprintln("File type: ", self.FileType))
	s.WriteString(fmt.Sprintln("set file attributes: \n ", self.Sattr3))

	return s.String()
}

type MknodReply struct {
	NFSStatus
	FileHandle []byte
	PostOpAttr
	Wcc
}

func (self *MknodReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	s.WriteString(fmt.Sprintln("File handle:", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintf("attributes:\n%v", self.PostOpAttr))
	s.WriteString(fmt.Sprintf("wcc:\n%v", self.Wcc))

	return s.String()
}

type RemoveCall struct {
	DirOpArg
}

func (self *RemoveCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("Dir op arg:\n ", self.DirOpArg))

	return s.String()
}

type RemoveReply struct {
	NFSStatus
	Wcc
}

func (self *RemoveReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	s.WriteString(fmt.Sprintln("wcc:\n", self.Wcc.String()))

	return s.String()
}

type RmdirCall struct {
	DirOpArg
}

func (self *RmdirCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("Dir op arg:\n ", self.DirOpArg))

	return s.String()
}

type RmdirReply struct {
	NFSStatus
	Wcc
}

func (self *RmdirReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	s.WriteString(fmt.Sprintln("wcc:\n", self.Wcc))

	return s.String()
}

type RenameCall struct {
	From DirOpArg
	To   DirOpArg
}

func (self *RenameCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("from:\n ", self.From))
	s.WriteString(fmt.Sprintln("to:\n ", self.To))

	return s.String()
}

type RenameReply struct {
	NFSStatus
	From Wcc
	To   Wcc
}

func (self *RenameReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	s.WriteString(fmt.Sprintln("from wcc:\n", self.From.String()))
	s.WriteString(fmt.Sprintln("to wcc:\n", self.To.String()))

	return s.String()
}

type LinkCall struct {
	FileHandle []byte
	DirOpArg
}

func (self *LinkCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("dir op arg:\n ", self.DirOpArg))

	return s.String()
}

type LinkReply struct {
	NFSStatus
	ObjAttr PostOpAttr
	Wcc
}

func (self *LinkReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.ObjAttr.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.ObjAttr.String()))
	}

	s.WriteString(fmt.Sprintln("wcc:\n", self.Wcc))

	return s.String()
}

type ReadDirCall struct {
	Handle   []byte
	Cookie   uint64
	Verifier uint64
	Count    uint32
}

func (self *ReadDirCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("handle: ", hex.EncodeToString(self.Handle)))
	s.WriteString(fmt.Sprintln("cookie: ", self.Cookie))
	s.WriteString(fmt.Sprintln("verifier: ", self.Verifier))
	s.WriteString(fmt.Sprintln("count: ", self.Count))

	return s.String()
}

type Entry struct {
	Field    uint64
	Filename []byte
	Cookie   uint64
}

type Entry3 struct {
	IsSet bool  `xdr:"union"`
	Entry Entry `xdr:"unioncase=1"`
}

type ReadDirArg struct {
	NFSStatus
	DirAttr  PostOpAttr
	Verifier uint64
}

type ReadDirReply struct {
	ReadDirArg
	Entries []*Entry
}

func (self *ReadDirReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.DirAttr.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.DirAttr.String()))
	}
	s.WriteString(fmt.Sprintln("verifier:", self.Verifier))
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

type ReadDirPlusCall struct {
	Handle   []byte
	Cookie   uint64
	Verifier uint64
	DirCount uint32
	MaxCount uint32
}

func (self *ReadDirPlusCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.Handle)))
	s.WriteString(fmt.Sprintln("cookie:", self.Cookie))
	s.WriteString(fmt.Sprintln("verifier:", self.Verifier))
	s.WriteString(fmt.Sprintln("dir count:", self.DirCount))
	s.WriteString(fmt.Sprintln("max count:", self.MaxCount))

	return s.String()
}

type EntryPlus struct {
	FileId   uint64
	FileName []byte
	Cookie   uint64
	Attr     PostOpAttr
	Handle   PostOpFH3
	// NextEntry *EntryPlus
}

func (e EntryPlus) FileHandleHex() string {
	return hex.EncodeToString(e.Handle.FileHandle)
}

func (e EntryPlus) FileNameStr() string {
	encoding := charset.GetStrCoding(e.FileName)
	switch encoding {
	case charset.GBK:
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(e.FileName) //将gbk再转换为utf-8
		return string(utf8Data)
	}
	return string(e.FileName)
}

type EntryPlus3 struct {
	IsSet bool      `xdr:"union"`
	Entry EntryPlus `xdr:"unioncase=1"`
}

type ReaDirPlusArg struct {
	NFSStatus
	DirAttrs   PostOpAttr
	CookieVerf uint64
}

type ReadDirPlusReply struct {
	ReaDirPlusArg
	Entries []*EntryPlus
}

func (self *ReadDirPlusReply) String() string {
	var s strings.Builder

	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.DirAttrs.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.DirAttrs.String()))
	}
	s.WriteString(fmt.Sprintln("cookie verifier:", self.CookieVerf))
	s.WriteString(fmt.Sprintln("entrys:"))
	for _, v := range self.Entries {
		s.WriteString(fmt.Sprintln("field:", v.FileId))
		fmt.Println(charset.GetStrCoding(v.FileName))
		utf8Data, _ := simplifiedchinese.GBK.NewDecoder().Bytes(v.FileName) //将gbk再转换为utf-8
		s.WriteString(fmt.Sprintln("    file name:", string(utf8Data)))
		s.WriteString(fmt.Sprintln("cookie:", v.Cookie))
		if v.Attr.HasSet() {
			s.WriteString(fmt.Sprintln("Attr:\n", v.Attr.String()))
		}
		s.WriteString(fmt.Sprintln("file handle:", v.Handle.String()))

	}

	return s.String()
}

type FsStatCall struct {
	FileHandle []byte
}

func (self *FsStatCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))

	return s.String()
}

type FsStatReply struct {
	NFSStatus
	ObjAttr        PostOpAttr
	TotalSize      uint64
	FreeSize       uint64
	AvailableSize  uint64
	TotalFiles     uint64
	FreeFiles      uint64
	AvailableFiles uint64
	// CacheHint is called "invarsec" in the nfs standard
	Invarsec uint32
}

func (self FsStatReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.ObjAttr.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.ObjAttr.String()))
	}
	s.WriteString(fmt.Sprintln("	total size:", self.TotalSize))
	s.WriteString(fmt.Sprintln("	free size:", self.FreeSize))
	s.WriteString(fmt.Sprintln("	available size:", self.AvailableSize))
	s.WriteString(fmt.Sprintln("	total files:", self.TotalFiles))
	s.WriteString(fmt.Sprintln("	free files:", self.FreeFiles))
	s.WriteString(fmt.Sprintln("	available files:", self.AvailableFiles))
	s.WriteString(fmt.Sprintln("	invarsec:", self.Invarsec))

	return s.String()
}

type FsInfoCall struct {
	FileHandle []byte
}

func (self *FsInfoCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	return s.String()
}

type FsInfoReply struct {
	NFSStatus
	Attr       PostOpAttr
	Rtmax      uint32
	Rtpref     uint32
	Rtmult     uint32
	Wtmax      uint32
	Wtpref     uint32
	Wtmult     uint32
	Dtpref     uint32
	Size       uint64
	TimeDelta  NFS3Time
	Properties uint32 //todo:format
}

func (self *FsInfoReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.Attr.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.Attr.String()))
	}
	s.WriteString(fmt.Sprintln("rtmax:", self.Rtmax))
	s.WriteString(fmt.Sprintln("rtpref:", self.Rtpref))
	s.WriteString(fmt.Sprintln("rtmult:", self.Rtmult))
	s.WriteString(fmt.Sprintln("wtmax:", self.Wtmax))
	s.WriteString(fmt.Sprintln("wtpref:", self.Wtpref))
	s.WriteString(fmt.Sprintln("dtpref:", self.Dtpref))
	s.WriteString(fmt.Sprintln("max file size:", self.Size))
	s.WriteString(fmt.Sprintln("time delta:", self.TimeDelta))
	s.WriteString(fmt.Sprintln("properties:", self.Properties))

	return s.String()
}

type PathConfCall struct {
	RootHandle []byte
}

func (self *PathConfCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.RootHandle)))
	return s.String()
}

type PathConfReply struct {
	NFSStatus
	ObjAttr         PostOpAttr
	LinkMax         uint32
	NameMax         uint32
	NoTrunc         uint32 //todo:bool
	ChownRestricted uint32 //todo:bool
	CaseInsensitive uint32 //todo:bool
	CasePreserving  uint32 //todo:bool
}

func (self *PathConfReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	if self.ObjAttr.HasSet() {
		s.WriteString(fmt.Sprintln("obj attributes:\n", self.ObjAttr.String()))
	}

	s.WriteString(fmt.Sprintln("link max:", self.LinkMax))
	s.WriteString(fmt.Sprintln("name max:", self.NameMax))
	s.WriteString(fmt.Sprintln("no trunc:", self.NoTrunc))
	s.WriteString(fmt.Sprintln("chown restricted:", self.ChownRestricted))
	s.WriteString(fmt.Sprintln("case insensitive:", self.CaseInsensitive))
	s.WriteString(fmt.Sprintln("case preserving:", self.CasePreserving))

	return s.String()
}

type CommitCall struct {
	FileHandle []byte
	Offset     uint64
	Count      uint32
}

func (self *CommitCall) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("file handle: ", hex.EncodeToString(self.FileHandle)))
	s.WriteString(fmt.Sprintln("offset: ", self.Offset))
	s.WriteString(fmt.Sprintln("count: ", self.Count))

	return s.String()
}

type CommitReply struct {
	NFSStatus
	Wcc
	Verifier uint64
}

func (self *CommitReply) String() string {
	var s strings.Builder
	s.WriteString(fmt.Sprintln("status:", self.NFSStatus))
	s.WriteString(fmt.Sprintln("wcc:\n", self.Wcc.String()))
	s.WriteString(fmt.Sprintln("verifier:", self.Verifier))

	return s.String()
}
