package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	managers "github.com/Alfazal007/gather-town/Managers"
	"github.com/Alfazal007/gather-town/types"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func wsHandler(w http.ResponseWriter, r *http.Request, wsManager *managers.WebSocketManager) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade the websocket connection")
		return
	}
	defer func() {
		wsManager.CleanUp(conn)
		conn.Close()
	}()
	for {
		t, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		messageIsOfCorrectType, typeOfMessage := wsManager.TypeChecker(message)
		if !messageIsOfCorrectType {
			continue
		}
		var messageInJsonFormat types.Message
		err = json.Unmarshal(message, &messageInJsonFormat)
		if err != nil {
			continue
		}
		if typeOfMessage == string(types.Conect) {
			var connectMessage types.ConectMessageSent
			_ = json.Unmarshal(messageInJsonFormat.Message, &connectMessage)
			canBeConnected := wsManager.RegisterUser(connectMessage, messageInJsonFormat.Username, messageInJsonFormat.Room)
			if canBeConnected {
				wsManager.ConnectMessageHandler(messageInJsonFormat, conn)
			}
		} else if typeOfMessage == string(types.Disconnect) {
			wsManager.DisconnectMessageHandler(messageInJsonFormat, conn, t)
		} else if typeOfMessage == string(types.TextMessage) {
			wsManager.SendTextMessage(messageInJsonFormat, conn, t)
		} else if typeOfMessage == string(types.PositionMessage) {
			wsManager.SendPositionMessage(messageInJsonFormat, conn, t)
		} else if typeOfMessage == string(types.InitiateCallRequest) {
			wsManager.HandleInitiateCallMessage(messageInJsonFormat, conn, t)
		} else if typeOfMessage == string(types.AcceptCallResponse) {
			wsManager.HandleAcceptCallMessage(messageInJsonFormat, conn, t)
		}
	}
}

func main() {
	wsManager := managers.WebSocketManager{
		Mutex:          sync.RWMutex{},
		RoomWithPeople: make(map[string]map[string]*websocket.Conn),
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, &wsManager)
	})
	err := http.ListenAndServe(":8001", nil)
	fmt.Println(err)
}
