package main

import (
	"log"
	"github.com/majalcmaj/tictactoe-server/server"
)

// func runWs(w http.ResponseWriter, r *http.Request) {
// 	var upgrader = websocket.Upgrader{}
// 	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
// 	c, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("Upgrade: ", err)
// 		return
// 	}

// 	defer c.Close()

// 	var msg Message
// 	for {
// 		mt, message, err := c.ReadMessage()
// 		if err != nil {
// 			log.Println("Read: ", err)
// 			break
// 		}

// 		json.Unmarshal([]byte(message), &msg)

// 		log.Printf("Index: %d", msg.Idx)
// 		err = c.WriteMessage(mt, message)
// 		if err != nil {
// 			log.Println("Write: ", err)
// 			break
// 		}
// 	}

// }

// func main() {
// 	var addr = flag.String("addr", "0.0.0.0:8080", "http service address")
// 	http.HandleFunc("/", runWs)
// 	log.Fatal(http.ListenAndServe(*addr, nil))
// }

func main() {
	server := server.NewServer(server.NewGamesService())
	log.Fatal(server.ServeForever(":8080"))
}
