package logic

import (
	"encoding/hex"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/config"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs/nfsv2"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/decoder/nfs/nfsv3"
	"github.com/funbinary/go_example/example/internal/proc_netaudit/hook"
	"github.com/funbinary/go_example/pkg/bfile"
	log "github.com/funbinary/go_example/pkg/blog"
	"github.com/funbinary/go_example/pkg/encoding/bjson"
	"github.com/funbinary/go_example/pkg/errors"
	"github.com/funbinary/go_example/pkg/rpc"
	"go.uber.org/atomic"
	"os"
	"sync"
	"time"
)

type NfsSession struct {
	cacheJsonFile string
	fhToName      sync.Map //map[string]string
	readSession   sync.Map //map[string]*NfsReadSession
	writeSession  sync.Map //map[string]*NfsWriteSession
}

func NewNfsSession(streamid string) *NfsSession {
	cacheJsonFile := bfile.Join(config.Config.CacheNetaudit, streamid+".json")
	fhCache := bfile.GetContents(cacheJsonFile)
	s := &NfsSession{
		cacheJsonFile: cacheJsonFile,
		//fhToName:      make(map[string]string),
		//readSession:   make(map[string]*NfsReadSession),
		//writeSession:  make(map[string]*NfsWriteSession),
	}
	if fhCache != "" {
		bjson.UnMarshal(fhCache, s.fhToName)
	}
	return s
}

func (s *NfsSession) Handle(app hook.Protocol) {

	if nfsapp, ok := app.(*nfs.Nfs); ok {
		switch nfsapp.CurMsgType {
		case rpc.RpcCall: //不进行处理，因为reply有带call信息
		case rpc.RpcReply:
			s.handleReply(nfsapp)
		}
	}
}

func (s *NfsSession) handleReply(app *nfs.Nfs) {
	switch app.Program {
	case rpc.NfsServiceID:
		switch rpc.NfsVersion(app.ProgVers) {
		case rpc.NfsV2:
			log.Infof("hook program:%v  version:%v produce:%v", app.Program, app.ProgVers, nfsv2.Procedure(app.Procedure))

			s.handleNfsV2(app)
		case rpc.NfsV3:
			log.Infof("hook program:%v  version:%v produce:%v", app.Program, app.ProgVers, nfsv3.Procedure(app.Procedure))
			s.handleNfsV3(app)
		}
	}

}

func (s *NfsSession) handleNfsV2(app *nfs.Nfs) {
	if nfsv2.Nfs2Status(app.RepHeader.State) != nfsv2.NFS2StatusOK {
		return
	}
	if app.Call == nil || app.Reply == nil {
		return
	}

	switch nfsv2.Procedure(app.Procedure) {

	case nfsv2.NFS2ProcedureLookUp:
		call := app.Call.(*nfsv2.LookUpCall)
		reply := app.Reply.(*nfsv2.LookUpReply)
		s.fhToName.Store(hex.EncodeToString(reply.FileHandle), call.FileNameStr())
		//s.fhToName[hex.EncodeToString(reply.FileHandle)] = call.FileNameStr()
		content := bjson.Marshal(s.fhToName)
		bfile.SetContents(s.cacheJsonFile, content)
	case nfsv2.NFS2ProcedureRead:
		call := app.Call.(*nfsv2.ReadCall)
		reply := app.Reply.(*nfsv2.ReadReply)
		//log.Infof("下载文件%s", call)
		filename := s.GetFileName(hex.EncodeToString(call.FileHandle))
		var readSession interface{}
		var ok = false
		if readSession, ok = s.readSession.Load(filename); !ok {
			readSession = NewNfs3ReadSession(filename, uint64(reply.Attr.Size))
			s.readSession.Store(filename, readSession)
		}
		err := readSession.(*NfsReadSession).OnReadV2(app, call, reply)
		if err != nil {
			log.Errorf("OnReadV2 error:%+v", err)
		}
	case nfsv2.NFS2ProcedureWrite:
		call := app.Call.(*nfsv2.WriteCall)
		reply := app.Reply.(*nfsv2.WriteReply)
		//log.Infof("上传文件%s", call)
		filename := s.GetFileName(hex.EncodeToString(call.FileHandle))
		var writeSession interface{}
		var ok = false
		if writeSession, ok = s.writeSession.Load(filename); !ok {
			writeSession = NewNfsWriteSession(filename)
			s.writeSession.Store(filename, writeSession)
		}
		if err := writeSession.(*NfsWriteSession).OnWriteV2(app, call, reply); err != nil {
			log.Errorf("%+v", err)
		}
	case nfsv2.NFS2ProcedureCreate:
		call := app.Call.(*nfsv2.CreateCall)
		reply := app.Reply.(*nfsv2.CreateReply)
		log.Infof("创建文件%s", call.DirOpArgs.FileNameStr())
		// 缓存句柄
		s.fhToName.Store(hex.EncodeToString(reply.FileHandle), call.FileNameStr())
		//s.fhToName[hex.EncodeToString(reply.FileHandle)] = call.DirOpArgs.FileNameStr()
		content := bjson.Marshal(s.fhToName)
		bfile.SetContents(s.cacheJsonFile, content)
	case nfsv2.NFS2ProcedureRemove:
		call := app.Call.(*nfsv2.RemoveCall)
		//reply := app.Reply.(*nfsv3.RemoveReply)
		log.Infof("删除文件%s", call.FileNameStr())
		//msg.SendToMain(&msg.NfsMsg{
		//	Protocol: msg.NfsProtocol,
		//	Version:  app.ProgVers,
		//	Operate:  msg.RemoveFile,
		//	FileName: call.FileNameStr(),
		//})
	case nfsv2.NFS2ProcedureRename:
		call := app.Call.(*nfsv2.RenameCall)
		log.Infof("重命名%s为%s", call.From.FileNameStr(), call.To.FileNameStr())
		//msg.SendToMain(&msg.NfsMsg{
		//	Protocol: msg.NfsProtocol,
		//	Version:  app.ProgVers,
		//	Operate:  msg.Rename,
		//	OldName:  call.From.FileNameStr(),
		//	NewName:  call.To.FileNameStr(),
		//})
	case nfsv2.NFS2ProcedureMkDir:
		call := app.Call.(*nfsv2.MkdirCall)
		//reply := app.Reply.(nfsv3.MkdirReply)
		log.Infof("创建目录%s", call.DirOpArgs.FileNameStr())
		//m := &msg.NfsMsg{
		//	Protocol: msg.NfsProtocol,
		//	Version:  app.ProgVers,
		//	Operate:  msg.CreateDir,
		//	FileName: call.DirOpArgs.FileNameStr(),
		//}
		//msg.SendToMain(m)
	case nfsv2.NFS2ProcedureRmDir:
		call := app.Call.(*nfsv2.RmdirCall)
		//reply := app.Reply.(*nfsv3.RmdirReply)
		log.Infof("删除目录%s", call.FileNameStr())
		//msg.SendToMain(&msg.NfsMsg{
		//	Protocol: msg.NfsProtocol,
		//	Version:  app.ProgVers,
		//	Operate:  msg.Rmdir,
		//	FileName: call.FileNameStr(),
		//})

	case nfsv2.NFS2ProcedureReadDir:

	case nfsv2.NFS2ProcedureNull:
	case nfsv2.NFS2ProcedureGetAttr:
	case nfsv2.NFS2ProcedureSetAttr:
	case nfsv2.NFS2ProcedureRoot:
	case nfsv2.NFS2ProcedureReadlink:
	case nfsv2.NFS2ProcedureWriteCache:
	case nfsv2.NFS2ProcedureLink:
	case nfsv2.NFS2ProcedureSymlink:
	case nfsv2.NFS2ProcedureStatFs:
	}

}

func (s *NfsSession) handleNfsV3(app *nfs.Nfs) {
	if nfsv3.NFSStatus(app.RepHeader.State) != nfsv3.NFSStatusOk {
		return
	}
	if app.Call == nil || app.Reply == nil {
		return
	}

	switch nfsv3.Procedure(app.Procedure) {
	case nfsv3.NFS3ProcedureMkDir:
		call := app.Call.(*nfsv3.MkdirCall)
		//reply := app.Reply.(nfsv3.MkdirReply)
		log.Infof("创建目录%s", call.Where.FileName())
		//m := &msg.NfsMsg{
		//	Protocol: msg.NfsProtocol,
		//	Version:  app.ProgVers,
		//	Operate:  msg.CreateDir,
		//	FileName: call.Where.FileName(),
		//}
		//msg.SendToMain(m)

	case nfsv3.NFS3ProcedureReadDirPlus:
		//call := app.Call.(*nfsv3.ReadDirPlusCall)
		reply := app.Reply.(*nfsv3.ReadDirPlusReply)
		for _, entry := range reply.Entries {
			s.fhToName.Store(entry.FileHandleHex(), entry.FileNameStr())
			//s.fhToName[entry.FileHandleHex()] = entry.FileNameStr()
		}
		content := bjson.Marshal(s.fhToName)
		bfile.SetContents(s.cacheJsonFile, content)
		//log.Infof("readdirplus")
	case nfsv3.NFS3ProcedureCreate:
		call := app.Call.(*nfsv3.CreateCall)
		reply := app.Reply.(*nfsv3.CreateReply)
		log.Infof("创建文件%s", call.FileName())
		// 缓存句柄
		s.fhToName.Store(hex.EncodeToString(reply.FileHandle.FileHandle), call.FileName())
		//s.fhToName[hex.EncodeToString(reply.FileHandle.FileHandle)] = call.FileName()

		content := bjson.Marshal(s.fhToName)
		bfile.SetContents(s.cacheJsonFile, content)

	case nfsv3.NFS3ProcedureRead:
		call := app.Call.(*nfsv3.ReadCall)
		reply := app.Reply.(*nfsv3.ReadReply)
		//log.Infof("下载文件%s", call)
		filename := s.GetFileName(hex.EncodeToString(call.FileHandle))
		var readSession interface{}
		var ok = false
		if readSession, ok = s.readSession.Load(filename); !ok {
			readSession = NewNfs3ReadSession(filename, reply.Attr.Attr.Filesize)
			s.readSession.Store(filename, readSession)
		}
		err := readSession.(*NfsReadSession).OnReadV3(app, call, reply)
		if err != nil {
			log.Errorf("OnReadV3 error:%+v", err)
		}
	case nfsv3.NFS3ProcedureRemove:
		call := app.Call.(*nfsv3.RemoveCall)
		//reply := app.Reply.(*nfsv3.RemoveReply)
		log.Infof("删除文件%s", call.FileName())
		//msg.SendToMain(&msg.NfsMsg{
		//	Protocol: msg.NfsProtocol,
		//	Version:  app.ProgVers,
		//	Operate:  msg.RemoveFile,
		//	FileName: call.FileName(),
		//})

	case nfsv3.NFS3ProcedureRmDir:
		call := app.Call.(*nfsv3.RmdirCall)
		//reply := app.Reply.(*nfsv3.RmdirReply)
		log.Infof("删除目录%s", call.FileName())
		//msg.SendToMain(&msg.NfsMsg{
		//	Protocol: msg.NfsProtocol,
		//	Version:  app.ProgVers,
		//	Operate:  msg.Rmdir,
		//	FileName: call.FileName(),
		//})

	case nfsv3.NFS3ProcedureRename:
		call := app.Call.(*nfsv3.RenameCall)
		log.Infof("重命名%s为%s", call.From.FileName(), call.To.FileName())
		//msg.SendToMain(&msg.NfsMsg{
		//	Protocol: msg.NfsProtocol,
		//	Version:  app.ProgVers,
		//	Operate:  msg.Rename,
		//	OldName:  call.From.FileName(),
		//	NewName:  call.To.FileName(),
		//})

	case nfsv3.NFS3ProcedureWrite:
		call := app.Call.(*nfsv3.WriteCall)
		reply := app.Reply.(*nfsv3.WriteReply)
		//log.Infof("上传文件%s", call)
		filename := s.GetFileName(hex.EncodeToString(call.Handle))
		var writeSession interface{}
		var ok = false
		if writeSession, ok = s.writeSession.Load(filename); !ok {
			writeSession = NewNfsWriteSession(filename)
			s.writeSession.Store(filename, writeSession)
		}
		if err := writeSession.(*NfsWriteSession).OnWriteV3(app, call, reply); err != nil {
			log.Errorf("%+v", err)
		}
	case nfsv3.NFS3ProcedureCommit:
		call := app.Call.(*nfsv3.CommitCall)
		//reply := app.Reply.(*nfsv3.CommitReply)
		//log.Infof("上传文件结束%s", call)
		filename := s.GetFileName(hex.EncodeToString(call.FileHandle))
		var writeSession interface{}
		var ok = false
		if writeSession, ok = s.writeSession.Load(filename); !ok {
			writeSession = NewNfsWriteSession(filename)
			s.writeSession.Store(filename, writeSession)
		}
		writeSession.(*NfsWriteSession).OnCommit(rpc.NfsV3)
	case nfsv3.NFS3ProcedureNull:
	case nfsv3.NFS3ProcedureGetAttr:
	case nfsv3.NFS3ProcedureSetAttr:
	case nfsv3.NFS3ProcedureLookup:
	case nfsv3.NFS3ProcedureAccess:
	case nfsv3.NFS3ProcedureReadlink:
	case nfsv3.NFS3ProcedureSymlink:
	case nfsv3.NFS3ProcedureMkNod:
	case nfsv3.NFS3ProcedureReadDir:
	case nfsv3.NFS3ProcedureLink:
	case nfsv3.NFS3ProcedureFSStat:
	case nfsv3.NFS3ProcedureFSInfo:
	case nfsv3.NFS3ProcedurePathConf:

	}

}

func (s *NfsSession) GetFileName(fh string) string {
	if v, ok := s.fhToName.Load(fh); ok {
		return v.(string)
	}
	return fh

}

type NfsReadSession struct {
	filename   string
	filepath   string
	fp         *os.File
	fileSize   uint64
	readySize  uint64
	ticker     *time.Ticker
	isComplete *atomic.Bool
}

func NewNfs3ReadSession(filename string, fileSize uint64) *NfsReadSession {
	filepath := bfile.Join(config.Config.CacheNetaudit, filename)
	fp, err := bfile.OpenFile(filepath, bfile.O_CREATE|bfile.O_WRONLY, bfile.DefaultPermOpen)
	if err != nil {
		log.Errorf("创建文件%s失败:%+v", filepath, err)
	}

	s := &NfsReadSession{
		filename:   filename,
		filepath:   filepath,
		fileSize:   fileSize,
		fp:         fp,
		ticker:     time.NewTicker(5 * time.Second),
		isComplete: atomic.NewBool(false),
	}
	go func() {

		select {
		case <-s.ticker.C:
			s.Completed()
		}

	}()
	return s
}

func (s *NfsReadSession) OnReadV2(app *nfs.Nfs, call *nfsv2.ReadCall, reply *nfsv2.ReadReply) (err error) {

	if s.fp == nil {
		s.fp, err = bfile.OpenFile(s.filepath, bfile.O_CREATE|bfile.O_WRONLY, bfile.DefaultPermOpen)
		if err != nil {
			return errors.Wrapf(err, "创建文件%s失败", s.filepath)
		}
	}
	_, err = s.fp.WriteAt(reply.Data, int64(call.Offset))
	if err != nil {
		return errors.Wrapf(err, "写入文件%s失败", s.filepath)
	}
	s.readySize += uint64(len(reply.Data))
	s.ticker.Reset(5 * time.Second)

	if s.readySize == s.fileSize {
		s.Completed()
	}
	return nil
}

func (s *NfsReadSession) OnReadV3(app *nfs.Nfs, call *nfsv3.ReadCall, reply *nfsv3.ReadReply) (err error) {

	if s.fp == nil {
		s.fp, err = bfile.OpenFile(s.filepath, bfile.O_CREATE|bfile.O_WRONLY, bfile.DefaultPermOpen)
		if err != nil {
			return errors.Wrapf(err, "创建文件%s失败", s.filepath)
		}
	}
	_, err = s.fp.WriteAt(reply.Data, int64(call.Offset))
	if err != nil {
		return errors.Wrapf(err, "写入文件%s失败", s.filepath)
	}
	s.readySize += uint64(len(reply.Data))
	s.ticker.Reset(5 * time.Second)

	if s.readySize == s.fileSize {
		s.Completed()
	}
	return nil
}

func (s *NfsReadSession) Completed() {
	if s.isComplete.Load() {
		return
	}
	s.isComplete.Store(true)
	stat, _ := s.fp.Stat()
	log.Infof("文件%s下载结束,大小:%d", s.filepath, stat.Size())
	// 写入完成
	//msg.SendToMain(&msg.NfsMsg{
	//	Protocol: msg.NfsProtocol,
	//	Version:  uint32(rpc.NfsV3),
	//	Operate:  msg.Download,
	//	FileName: s.filename,
	//})
}

type NfsWriteSession struct {
	filename   string
	filepath   string
	fp         *os.File
	fileSize   uint64
	isComplete *atomic.Bool
	ticker     *time.Ticker
	writeSize  uint64
}

func NewNfsWriteSession(filename string) *NfsWriteSession {
	filepath := bfile.Join(config.Config.CacheNetaudit, filename)
	fp, err := bfile.OpenFile(filepath, bfile.O_CREATE|bfile.O_WRONLY, bfile.DefaultPermOpen)
	if err != nil {
		log.Errorf("创建文件%s失败:%+v", filepath, err)
	}

	s := &NfsWriteSession{
		filename:   filename,
		filepath:   filepath,
		fp:         fp,
		isComplete: atomic.NewBool(false),
		ticker:     time.NewTicker(5 * time.Second),
		writeSize:  0,
	}
	go func() {

		select {
		case <-s.ticker.C:
			s.OnCommit(rpc.NfsV3)
		}

	}()
	return s
}

func (s *NfsWriteSession) OnWriteV2(app *nfs.Nfs, call *nfsv2.WriteCall, reply *nfsv2.WriteReply) (err error) {

	if s.fp == nil {
		s.fp, err = bfile.OpenFile(s.filepath, bfile.O_CREATE|bfile.O_WRONLY, bfile.DefaultPermOpen)
		if err != nil {
			return errors.Wrapf(err, "创建文件%s失败", s.filepath)
		}
	}
	var wLen int
	wLen, err = s.fp.WriteAt(call.Data, int64(call.Offset))
	log.Infof("写入文件%s 写入大小:%d 元数据大小:%d", s.filename, wLen, len(call.Data))
	s.writeSize += uint64(len(call.Data))
	s.ticker.Reset(5 * time.Second)
	if err != nil {
		return errors.Wrapf(err, "写入文件%s失败", s.filepath)
	}
	return nil
}

func (s *NfsWriteSession) OnWriteV3(app *nfs.Nfs, call *nfsv3.WriteCall, reply *nfsv3.WriteReply) (err error) {

	if s.fp == nil {
		s.fp, err = bfile.OpenFile(s.filepath, bfile.O_CREATE|bfile.O_WRONLY, bfile.DefaultPermOpen)
		if err != nil {
			return errors.Wrapf(err, "创建文件%s失败", s.filepath)
		}
	}
	var wLen int
	wLen, err = s.fp.WriteAt(call.Data, int64(call.Offset))
	log.Infof("写入文件%s 写入大小:%d 元数据大小:%d", s.filename, wLen, len(call.Data))
	s.writeSize += uint64(len(call.Data))
	s.ticker.Reset(5 * time.Second)
	if err != nil {
		return errors.Wrapf(err, "写入文件%s失败", s.filepath)
	}
	return nil
}

func (s *NfsWriteSession) OnCommit(v rpc.NfsVersion) {

	if s.isComplete.Load() {
		return
	}
	s.isComplete.Store(true)
	stat, _ := s.fp.Stat()
	log.Infof("文件%s上传结束,大小:%d", s.filepath, stat.Size())
	s.fp.Sync()
	s.fp.Close()

	// 写入完成
	//msg.SendToMain(&msg.NfsMsg{
	//	Protocol: msg.NfsProtocol,
	//	Version:  uint32(v),
	//	Operate:  msg.Upload,
	//	FileName: s.filename,
	//})

}
