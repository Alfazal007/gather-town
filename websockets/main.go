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
		_ = conn.Close()
	}()
	for {
		fmt.Println("town")
		wsManager.Mutex.RLock()
		fmt.Println(wsManager.RoomWithPeople)
		wsManager.Mutex.RUnlock()
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

func vsHandler(w http.ResponseWriter, r *http.Request, vsManager *managers.VideoRoomManager) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Failed to upgrade the websocket connection")
		return
	}
	defer func() {
		vsManager.CleanUp(conn)
		_ = conn.Close()
	}()
	for {
		fmt.Println("Vudei")
		vsManager.Mutex.RLock()
		fmt.Println(vsManager.RoomWithTwoPeople)
		vsManager.Mutex.RUnlock()
		t, message, err := conn.ReadMessage()
		if err != nil {
			break
		}
		messageIsOfCorrectType, typeOfMessage := vsManager.TypeChecker(message)
		if !messageIsOfCorrectType {
			continue
		}
		var messageInJsonFormat types.VideoMessage
		err = json.Unmarshal(message, &messageInJsonFormat)
		if err != nil {
			continue
		}
		if typeOfMessage == string(types.CreateRoomMessage) {
			var createRoomMessage types.CreateRoom
			_ = json.Unmarshal(messageInJsonFormat.Message, &createRoomMessage)
			canBeConnected := vsManager.RegisterUserForVideo(createRoomMessage.Token, messageInJsonFormat.Username)
			if canBeConnected {
				success := vsManager.CreateRoomForVideo(messageInJsonFormat, conn, t)
				if !success {
					return
				}
			}
		} else if typeOfMessage == string(types.IceCandidateMessage) {
			vsManager.ForwardIceCandidates(messageInJsonFormat, conn, t)
		} else if typeOfMessage == string(types.SDPRoomMessage) {
			vsManager.ForwardSDPData(messageInJsonFormat, conn, t)
		} else if typeOfMessage == string(types.JoinRoomMessage) {
			var connectMessage types.JoinRoom
			_ = json.Unmarshal(messageInJsonFormat.Message, &connectMessage)
			canBeConnected := vsManager.RegisterUserForVideo(connectMessage.Token, messageInJsonFormat.Username)
			if canBeConnected {
				vsManager.JoinRoomBySecondPerson(messageInJsonFormat, conn, t)
			}
		} else if typeOfMessage == string(types.DisconnectMessage) {
			vsManager.DisconnectVideoCall(messageInJsonFormat, conn, t)
		}
	}
}

func main() {
	wsManager := managers.WebSocketManager{
		Mutex:          sync.RWMutex{},
		RoomWithPeople: make(map[string]map[string]*websocket.Conn),
	}
	vsManager := managers.VideoRoomManager{
		RoomWithTwoPeople: make(map[string]managers.TwoPeople),
		Mutex:             sync.RWMutex{},
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(w, r, &wsManager)
	})
	http.HandleFunc("/video", func(w http.ResponseWriter, r *http.Request) {
		vsHandler(w, r, &vsManager)
	})
	err := http.ListenAndServe("0.0.0.0:8001", nil)
	//	err := http.ListenAndServe(":8001", nil)
	fmt.Println(err)
}
