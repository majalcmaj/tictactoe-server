package server

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	// "github.com/majalcmaj/tictactoe-server/game"
)

const sleepBeforeGameRemovalSeconds = 5
const maxConnections = 20
const maxGames = 3
const playersPerGame = 2

type IGamesService interface {
	RegisterGame(password string) (int, error)
	JoinGame(idx int, nick string, conn *Connection) error
	RemoveGame(idx int)
}

type GamesService struct {
	wsUpgrader          websocket.Upgrader
	startedGamesCounter int
	currentGames        map[int]*gameTransport
	accessLock          sync.Mutex
}

func NewGamesService() *GamesService {
	result := GamesService{
		wsUpgrader:          websocket.Upgrader{},
		startedGamesCounter: 0,
		currentGames:        make(map[int]*gameTransport),
	}
	result.wsUpgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return &result
}

func (gamesService *GamesService) RegisterGame(password string) (int, error) {
	if len(gamesService.currentGames) < maxGames {
		gameIdx := gamesService.nextGameIdx()
		game := newGameTransport(password, gamesService.getGameFinishedCallback(gameIdx))
		gamesService.currentGames[gameIdx] = game
		go gamesService.removeGameAfterTimeout(gameIdx)
		return gameIdx, nil
	} else {
		return 0, errors.New("maximum games count exceeded")
	}
}

func (gamesService *GamesService) JoinGame(idx int, password string, conn *Connection) error {
	if game := gamesService.currentGames[idx]; game != nil {
		game.accessLock.Lock()
		defer game.accessLock.Unlock()
		if game.checkPassword(password) {
			if game.canAddPlayer() {
				if err := gamesService.createPlayerForGame(game, conn); err == nil {
					log.Println("Player connected to game", idx)
					started := game.startIfReady()
					log.Println("Starting game: ", started)
					return nil
				} else {
					return err
				}
			}
			return fmt.Errorf("Game already has two players %d ", idx)
		}
		return fmt.Errorf("Wrong password for game %d ", idx)
	}
	return fmt.Errorf("No game with index %d exists", idx)
}

func (gamesService *GamesService) RemoveGame(idx int) {
	delete(gamesService.currentGames, idx)
}

func (gamesService *GamesService) upgradeToWebsocket(conn *Connection) (*websocket.Conn, error) {
	return gamesService.wsUpgrader.Upgrade(conn.w, conn.r, nil)
}

func (gamesService *GamesService) nextGameIdx() int {
	gamesService.accessLock.Lock()
	defer gamesService.accessLock.Unlock()
	result := gamesService.startedGamesCounter
	gamesService.startedGamesCounter++
	return result
}

func (gameService *GamesService) createPlayerForGame(game *gameTransport, conn *Connection) error {
	if ws, err := gameService.upgradeToWebsocket(conn); err == nil {
		player := playerTransport{websocket: ws}
		game.addPlayer(&player)
		return nil
	} else {
		return err
	}
}

func (gamesService *GamesService) removeGameAfterTimeout(idx int) {
	time.Sleep(time.Second * sleepBeforeGameRemovalSeconds)
	game := gamesService.currentGames[idx]
	game.accessLock.Lock()
	defer game.accessLock.Unlock()
	if game.connectedPlayersCount() == 0 {
		log.Printf("Removing game with index %d because no player has connected within timeout.", idx)
		gamesService.RemoveGame(idx)
	}
}

func (gamesService *GamesService) getGameFinishedCallback(index int) func() {
	return func() {
		gamesService.accessLock.Lock()
		defer gamesService.accessLock.Unlock()
		delete(gamesService.currentGames, index)
	}
}
