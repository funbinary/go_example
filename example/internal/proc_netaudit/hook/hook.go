package hook

type Hook interface {
	Hook(streamid string, app Protocol)
}

type Protocol interface {
	Protocol() string
}

func SetHook(h Hook) {
	hook = h
}

func Fire(streamid string, app Protocol) {
	hook.Hook(streamid, app)
}

var hook Hook = &defaultHook{}

type defaultHook struct{}

func (h *defaultHook) Hook(streamid string, app Protocol) {
	return
}
