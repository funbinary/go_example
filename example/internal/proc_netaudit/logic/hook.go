package logic

import (
	"github.com/funbinary/go_example/example/internal/proc_netaudit/hook"
	"sync"
)

type Hook struct {
	sessions sync.Map
}

func NewHook() *Hook {
	return &Hook{}
}

type Session interface {
	Handle(app hook.Protocol)
}

func (h *Hook) Hook(streamid string, app hook.Protocol) {
	if app == nil {
		return
	}
	var v interface{}
	var ok bool
	if v, ok = h.sessions.Load(streamid); !ok {
		switch app.Protocol() {
		case "nfs":
			v = NewNfsSession(streamid)
			h.sessions.Store(streamid, v)
		default:
			return
		}
	}
	session := v.(Session)
	session.Handle(app)
}
