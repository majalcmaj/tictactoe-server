package common

import (
	"log"
	"encoding/json"
	"fmt"
)

type EventType int

const (
	PlayerMove   EventType = 1
	PlayerTurn   EventType = 2
	MoveMade     EventType = 3
	GameFinished EventType = 4
)

type PlayerType int

const (
	None PlayerType = 0
	Tic  PlayerType = 1
	Tac  PlayerType = 2
)

type jsonEvent struct {
	EventType EventType       `json:"type"`
	Event     json.RawMessage `json:"event"`
}

type PlayerMoveEvent struct {
	Player   PlayerType `json:"player"`
	FieldIdx int        `json:"fieldIdx"`
}

func NewPlayerMoveEvent(player PlayerType, index int) *PlayerMoveEvent {
	return &PlayerMoveEvent{player, index}
}

type PlayerTurnEvent struct {
	Player PlayerType `json:"player"`
}

func NewPlayerTurnEvent(player PlayerType) *PlayerTurnEvent {
	return &PlayerTurnEvent{player}
}

type MoveMadeEvent struct {
	Player PlayerType `json:"player"`
	Index  int        `json:"fieldIdx"`
}

func NewMoveMadeEvent(player PlayerType, index int) *MoveMadeEvent {
	return &MoveMadeEvent{player, index}
}

type GameFinishedEvent struct {
	Player PlayerType `json:"player"`
}

func NewGameFinishedEvent(player PlayerType) *GameFinishedEvent {
	return &GameFinishedEvent{player}
}

func EncodeEvent(event interface{}) ([]byte, error) {
	if evType, ok := getEventType(event); ok {
		if bytes, err := json.Marshal(event); err == nil {
			return json.Marshal(&jsonEvent{evType, bytes})
		} else {
			return nil, err
		}
	}
	return nil, fmt.Errorf("Unknown event type: %s", event)
}

func getEventType(event interface{}) (EventType, bool) {
	switch event.(type) {
	case *PlayerMoveEvent:
		return PlayerMove, true
	case *PlayerTurnEvent:
		return PlayerTurn, true
	case *MoveMadeEvent:
		return MoveMade, true
	case *GameFinishedEvent:
		return GameFinished, true
	}
	return PlayerMove, false
}

func DecodeEvent(message []byte) (interface{}, error) {
	var event jsonEvent
	if err := json.Unmarshal(message, &event); err == nil {
		log.Println("Got event of type ", event.EventType)
		switch event.EventType {
		case PlayerMove:
			{
				var ev PlayerMoveEvent
				err := json.Unmarshal(event.Event, &ev)
				return ev, err
			}
		case PlayerTurn:
			{
				var ev PlayerTurnEvent
				err := json.Unmarshal(event.Event, &ev)
				return ev, err
			}
		case MoveMade:
			{
				var ev MoveMadeEvent
				err := json.Unmarshal(event.Event, &ev)
				return ev, err
			}
		case GameFinished:
			{
				var ev GameFinishedEvent
				err := json.Unmarshal(event.Event, &ev)
				return ev, err
			}
		default:
			return nil, fmt.Errorf("Cannot decode event of type %d", event.EventType)
		}
	} else {
		return nil, err
	}
}
