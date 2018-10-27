package test

import (
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/majalcmaj/tictactoe-server/transport"
	"github.com/stretchr/testify/require"
)

func expectResultWithin(t *testing.T, ms int, check func()) {
	passChan := make(chan int)
	go func() {
		check()
		passChan <- 0
	}()

	select {
	case <-passChan:
		return
	case <-time.NewTimer(time.Duration(ms) * time.Second).C:
		require.Fail(t, "No event received within timeout")
	}
}

func TestPlayersConnected(t *testing.T) {
	ss := transport.NewSocketServer()
	channel := ss.ControlEventsEmitter()
	go ss.NewConnection(&WebsocketMock{})
	go ss.NewConnection(&WebsocketMock{})

	expectResultWithin(t, 5, func() {
		msg := <-channel
		require.Equal(t, &transport.ControlEvent{transport.PeerConnected, 0}, msg)
	})
	expectResultWithin(t, 5, func() {
		msg := <-channel
		require.Equal(t, &transport.ControlEvent{transport.PeerConnected, 1}, msg)
	})

}

func TestMessageSending(t *testing.T) {
	req := require.New(t)
	ss := transport.NewSocketServer()

	bytes := [...]byte{1, 2, 3}
	ws := WebsocketMock{sendCallback: func(msgType int, data []byte) {
		req.Equal(msgType, websocket.TextMessage)
		req.Equal(bytes[:], data)
	}}
	ss.NewConnection(&ws)
	<-ss.ControlEventsEmitter()

	sendChan, _, err := ss.GetPeerChannels(0)
	req.Empty(err)
	require.Empty(t, err)
	expectResultWithin(t, 5, func() {
		sendChan <- bytes[:]
	})
}

func TestMessagesReceiving(t *testing.T) {
	req := require.New(t)
	ss := transport.NewSocketServer()

	bytes := [...]byte{1, 2, 3}
	ws := WebsocketMock{readType: websocket.TextMessage, readData: bytes[:]}
	ss.NewConnection(&ws)
	<-ss.ControlEventsEmitter()

	_, receiveChan, err := ss.GetPeerChannels(0)
	req.Empty(err)
	expectResultWithin(t, 5, func() {
		recvBytes := <-receiveChan
		req.Equal(bytes[:], recvBytes)
	})
}
