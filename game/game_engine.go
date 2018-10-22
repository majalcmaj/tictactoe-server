package game

import "github.com/majalcmaj/tictactoe-server/common"

const fieldsCount = 9

var winningMoves = [...]([3]int){
	[...]int{0, 1, 2}, [...]int{3, 4, 5}, [...]int{6, 7, 8},
	[...]int{0, 3, 6}, [...]int{1, 4, 7}, [...]int{2, 5, 8},
	[...]int{0, 4, 8}, [...]int{2, 4, 6},
}

type gameEventsHandler interface {
	playerTurn(common.PlayerType)
	playerMoved(common.PlayerType, int)
	gameFinished(common.PlayerType)
}

type gameEngine struct {
	gameState    [fieldsCount]common.PlayerType
	evHandler    gameEventsHandler
	isPlayerTic  bool
	movesCounter int
}

func newGameEngine(evHandler gameEventsHandler) *gameEngine {
	gameEngine := gameEngine{evHandler: evHandler, movesCounter: 0}
	for i := 0; i < fieldsCount; i++ {
		gameEngine.gameState[i] = common.None
	}
	gameEngine.evHandler.playerTurn(gameEngine.currentPlayer())
	return &gameEngine
}

func (gameEng *gameEngine) movePlayer(player common.PlayerType, index int) bool {
	if !gameEng.correctPlayerMoves(player) {
		return false
	}
	if !gameEng.correctFieldIdx(index) {
		return false
	}
	gameEng.makeMove(player, index)
	finished, winner := gameEng.checkForFinish()
	gameEng.evHandler.playerMoved(player, index)
	if finished {
		gameEng.evHandler.gameFinished(winner)
	} else {
		gameEng.movesCounter++
		gameEng.changePlayer()
		gameEng.evHandler.playerTurn(gameEng.currentPlayer())
	}
	return finished
}

func (gameEng *gameEngine) correctPlayerMoves(player common.PlayerType) bool {
	return gameEng.currentPlayer() == player
}

func (gameEng *gameEngine) correctFieldIdx(index int) bool {
	if index >= 0 && index < fieldsCount {
		return gameEng.gameState[index] == common.None
	}
	return false
}

func (gameEng *gameEngine) currentPlayer() common.PlayerType {
	if gameEng.isPlayerTic {
		return common.Tic
	}
	return common.Tac
}

func (gameEng *gameEngine) makeMove(player common.PlayerType, index int) {
	gameEng.gameState[index] = player
}

func (gameEng *gameEngine) checkForFinish() (bool, common.PlayerType) {
	if !(gameEng.movesCounter < fieldsCount) {
		return false, common.None
	}
	if gameEng.currentPlayerWon() {
		return true, gameEng.currentPlayer()
	}
	return false, common.None
}

func (gameEng *gameEngine) currentPlayerWon() bool {
	currentPlayer := gameEng.currentPlayer()
	playerMovesMap := common.Map(gameEng.gameState[:], func(f common.PlayerType) bool {
		return f == currentPlayer
	})
	gameEng.checkForWinningMoves(playerMovesMap)
	return false
}

func (gameEng *gameEngine) checkForWinningMoves(playerMovesMap []bool) bool {
	for _, move := range winningMoves {
		if playerMovesMap[move[0]] && playerMovesMap[move[1]] && playerMovesMap[move[2]] {
			return true
		}
	}
	return false
}

func (gameEng *gameEngine) changePlayer() {
	gameEng.isPlayerTic = !gameEng.isPlayerTic
}
