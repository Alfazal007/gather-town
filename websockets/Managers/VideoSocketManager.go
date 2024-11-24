package managers

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Person struct {
	Username string
	Conn     *websocket.Conn
}

type TwoPeople struct {
	Person1 Person
	Person2 Person
}

type VideoRoomManager struct {
	RoomWithTwoPeople map[string]TwoPeople
	Mutex             sync.RWMutex
}
