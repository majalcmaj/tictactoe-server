package server

import (
	"encoding/json" // "flag"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

type gameCreatedMessage struct {
	GameIndex int `json:"gameIndex"`
}

type RegistrationMessage struct {
	Password string `json:"password"`
}

type JoinMessage struct {
	Password string `json:"password"`
}

type Connection struct {
	w http.ResponseWriter
	r *http.Request
}

type GameInfo struct {
	Index int
}

type Server struct {
	wsUpgrader   websocket.Upgrader
	router       *mux.Router
	gamesService IGamesService
}

// NewServer creates a server instance
func NewServer(gameRegister IGamesService) *Server {
	router := mux.NewRouter()
	result := &Server{
		router:       router,
		gamesService: gameRegister,
	}
	result.setupRouter()
	return result
}

//ServeForever starts serving on given address "hostname:port"
func (server *Server) ServeForever(listenOn string) error {
	log.Println("Server listens on ", listenOn)
	return http.ListenAndServe(listenOn, server.router)
}

func (server *Server) createGame(w http.ResponseWriter, r *http.Request) {
	var registrationMessage RegistrationMessage
	handlePost(w, r, &registrationMessage, func() {
		if gameIndex, err := server.gamesService.RegisterGame(registrationMessage.Password); err == nil {
			log.Printf("Register game called, passsword: '%s'\n", registrationMessage.Password)
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(gameCreatedMessage{GameIndex: gameIndex})
		} else {
			log.Println("Failure: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func (server *Server) joinGame(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	// var joinGameMessage RegistrationMessage
	// handlePost(w, r, &joinGameMessage, func() {
	if idx, err := strconv.Atoi(params["game_id"]); err == nil {
		// if err := server.gamesService.JoinGame(idx, joinGameMessage.Password, &Connection{w, r}); err != nil {
		if err := server.gamesService.JoinGame(idx, "test", &Connection{w, r}); err != nil {
			log.Println("Failure: ", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	} else {
		log.Println("Failure: ", err)
		w.WriteHeader(http.StatusBadRequest)
	}
	// })
}

func handlePost(w http.ResponseWriter, r *http.Request, message interface{}, callback func()) {
	if err := json.NewDecoder(r.Body).Decode(&message); err == nil {
		callback()
	} else {
		log.Println("Got malformed JSON: ", err)
		w.WriteHeader(http.StatusBadRequest)
	}
}

func (server *Server) setupRouter() {
	server.router.HandleFunc("/game", server.createGame).Methods("POST")
	server.router.HandleFunc("/game/{game_id}/players", server.joinGame).Methods("GET")
}
