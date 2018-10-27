package test

type WebsocketMock struct {
	readType  int
	readData  []byte
	readError error

	sendError    error
	sendCallback func(int, []byte)
}

func (wsMock *WebsocketMock) ReadMessage() (messageType int, p []byte, err error) {
	return wsMock.readType, wsMock.readData, wsMock.readError
}

func (wsMock *WebsocketMock) WriteMessage(messageType int, data []byte) error {
	wsMock.sendCallback(messageType, data)
	return wsMock.sendError
}
