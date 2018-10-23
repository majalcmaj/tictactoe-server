package game

import "github.com/majalcmaj/tictactoe-server/common"

type GameAdapter struct {
	inEventChan  <-chan interface{}
	outEventChan chan<- interface{}
	gameEngine   *gameEngine
}

func NewGameAdapter(inEventChan <-chan interface{}, outEventChan chan<- interface{}) *GameAdapter {
	gameAdapter := GameAdapter{inEventChan, outEventChan, nil}
	gameAdapter.gameEngine = newGameEngine(&gameAdapter)
	return &gameAdapter
}

func (gameAdapter *GameAdapter) Run() {
	for event := range gameAdapter.inEventChan {
		if finished := gameAdapter.handleEvent(event); finished {
			break
		}
	}
	close(gameAdapter.outEventChan)
}

func (gameAdapter *GameAdapter) handleEvent(anyEvent interface{}) bool {
	switch event := anyEvent.(type) {
	case *common.MoveMadeEvent:
		return gameAdapter.gameEngine.movePlayer(event.Player, event.Index)
	}
	return false
}

func (gameAdapter *GameAdapter) playerTurn(player common.PlayerType) {
	gameAdapter.outEventChan <- common.NewPlayerTurnEvent(player)
}

func (gameAdapter *GameAdapter) playerMoved(player common.PlayerType, fieldIndex int) {
	gameAdapter.outEventChan <- common.NewMoveMadeEvent(player, fieldIndex)
}

func (gameAdapter *GameAdapter) gameFinished(player common.PlayerType) {
	gameAdapter.outEventChan <- common.NewGameFinishedEvent(player)
}
