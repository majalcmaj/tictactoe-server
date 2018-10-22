package server

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/majalcmaj/tictactoe-server/common"
	"github.com/majalcmaj/tictactoe-server/game"
)

type playerTransport struct {
	nick      string
	websocket *websocket.Conn
}

func (playerTrans *playerTransport) closeWs() {
	playerTrans.websocket.Close()
}

func (playerTrans *playerTransport) send(msg []byte) error {
	return playerTrans.websocket.WriteMessage(websocket.TextMessage, msg)
}

type gameTransport struct {
	players        []*playerTransport
	accessLock     sync.Mutex
	password       string
	gameEngine     *game.GameAdapter
	inEventChan    chan<- interface{}
	outEventChan   <-chan interface{}
	finishCallback func()
}

func newGameTransport(password string, finishCallback func()) *gameTransport {
	gameIns := gameTransport{
		players:        make([]*playerTransport, 0, playersPerGame),
		password:       password,
		finishCallback: finishCallback,
	}
	return &gameIns
}

func (gameTrans *gameTransport) connectedPlayersCount() int {
	return len(gameTrans.players)
}

func (gameTrans *gameTransport) canAddPlayer() bool {
	return len(gameTrans.players) < playersPerGame
}

func (gameTrans *gameTransport) addPlayer(player *playerTransport) {
	gameTrans.players = append(gameTrans.players, player)
}

func (gameTrans *gameTransport) checkPassword(password string) bool {
	return gameTrans.password == password
}

func (gameTrans *gameTransport) startIfReady() bool {
	if len(gameTrans.players) == playersPerGame {
		inEventChan := make(chan interface{})
		outEventChan := make(chan interface{})
		gameTrans.inEventChan = inEventChan
		gameTrans.outEventChan = outEventChan
		gameTrans.startSending()
		gameTrans.gameEngine = game.NewGameAdapter(inEventChan, outEventChan)
		go gameTrans.gameEngine.Run()
		gameTrans.startReceiving()
		return true
	}
	return false
}

func (gameTrans *gameTransport) startSending() {
	go func() {
		for {
			event, ok := <-gameTrans.outEventChan
			if !ok {
				gameTrans.cleanup()
				break
			}
			if encEvent, err := common.EncodeEvent(event); err == nil {
				log.Println("Sending event to players: ", encEvent)
				for _, player := range gameTrans.players {
					if err := player.send(encEvent); err != nil {
						log.Println("Error during message sending")
						gameTrans.cleanup()
						return
					}
				}
			} else {
				log.Println("Error occured during message: ", event, "encoding: ", err)
			}
		}
	}()
}

func (gameTrans *gameTransport) startReceiving() {
	go gameTrans.startReceivingForPlayer(0)
	go gameTrans.startReceivingForPlayer(1)
}

func (gameTrans *gameTransport) startReceivingForPlayer(playerIdx int) {
	playerTrans := gameTrans.players[playerIdx]
	for {
		_, message, err := playerTrans.websocket.ReadMessage()
		if err != nil {
			log.Println("Error during receiving form player", playerIdx, ":", err)
			return
		}
		if event, err := common.DecodeEvent(message); err == nil {
			log.Println("Got event from player: ", event)
			gameTrans.inEventChan <- event
		} else {
			log.Println("Error during message decoding:", err)
		}
	}
}

func (gameTrans *gameTransport) cleanup() {
	close(gameTrans.inEventChan)
	for _, playerTrans := range gameTrans.players {
		playerTrans.closeWs()
	}
	gameTrans.finishCallback()
}
