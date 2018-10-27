package transport

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type peer struct {
	ws       Websocket
	received chan []byte
	toSend   chan []byte
}

func (p *peer) startHandling() {
	go p.startReceiving()
	go p.startSending()
}

func (p *peer) startReceiving() {
	ok := true
	for ok {
		_, msg, error := p.ws.ReadMessage()
		if error != nil {
			close(p.received)
		}
		p.received <- msg
	}
}

func (p *peer) startSending() {
	for {
		msg, ok := <-p.toSend
		if !ok {
			return
		}
		if err := p.ws.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
	}
}

type SocketServer struct {
	peers                 map[int]*peer
	connectionsCounter    int
	controlEventsEmmitter chan *ControlEvent
	mutex                 sync.Mutex
}

func NewSocketServer() *SocketServer {
	return &SocketServer{peers: make(map[int]*peer), controlEventsEmmitter: make(chan *ControlEvent)}
}

func (socketServer *SocketServer) NewConnection(websocket Websocket) {
	socketServer.mutex.Lock()
	defer socketServer.mutex.Unlock()
	newPeer := &peer{websocket, make(chan []byte), make(chan []byte)}
	idx := socketServer.connectionsCounter
	socketServer.connectionsCounter++
	socketServer.peers[idx] = newPeer
	newPeer.startHandling()
	go socketServer.emitNewConnectionEvent(idx)
}

func (socketServer *SocketServer) ControlEventsEmitter() <-chan *ControlEvent {
	return socketServer.controlEventsEmmitter
}

func (socketServer *SocketServer) GetPeerChannels(idx int) (chan<- []byte, <-chan []byte, error) {
	peer := socketServer.peers[idx]
	if peer != nil {
		return peer.toSend, peer.received, nil
	}
	return nil, nil, fmt.Errorf("No peer with index %d", idx)
}

func (socketServer *SocketServer) emitNewConnectionEvent(idx int) {
	socketServer.controlEventsEmmitter <- &ControlEvent{PeerConnected, idx}
}
