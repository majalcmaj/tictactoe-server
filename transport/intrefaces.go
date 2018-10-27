package transport

type Websocket interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
}

type ControlEventType int

const (
	PeerConnected    ControlEventType = 1
	PeerDisconnected ControlEventType = 2
)

type ControlEvent struct {
	EventType ControlEventType
	PeerIndex int
}
