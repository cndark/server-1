package app

// ============================================================================

type socket_t interface {
	Close()
	Send([]byte)
	RemoteIP() string

	OnData(func(buf []byte))
	OnClose(func())
}
