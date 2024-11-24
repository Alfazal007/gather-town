package managers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Alfazal007/gather-town/types"
	"github.com/gorilla/websocket"
)

type WebSocketManager struct {
	RoomWithPeople map[string]map[string]*websocket.Conn
	Mutex          sync.RWMutex
}

func (wsManager *WebSocketManager) RegisterUser(message types.ConectMessageSent, username string, roomToJoin string) bool {
	tokenExtracted := message.Token
	if len(tokenExtracted) < 5 || len(roomToJoin) < 4 {
		return false
	}
	url := fmt.Sprintf("%v/%v/token/%v/username/%v", types.BackendUrl, roomToJoin, tokenExtracted, username)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (wsManager *WebSocketManager) ConnectMessageHandler(messageSent types.Message, connection *websocket.Conn) {
	wsManager.Mutex.Lock()
	_, roomExists := wsManager.RoomWithPeople[messageSent.Room]
	if roomExists {
		wsManager.RoomWithPeople[messageSent.Room][messageSent.Username] = connection
	} else {
		newMap := make(map[string]*websocket.Conn)
		newMap[messageSent.Username] = connection
		wsManager.RoomWithPeople[messageSent.Room] = newMap
	}
	wsManager.Mutex.Unlock()
}

func (wsManager *WebSocketManager) DisconnectMessageHandler(messageSent types.Message, connection *websocket.Conn, messageType int) {
	wsManager.Mutex.Lock()
	roomToBeBroadcastedIn, roomExists := wsManager.RoomWithPeople[messageSent.Room]
	if roomExists {
		connectionToBeClosed := wsManager.RoomWithPeople[messageSent.Room][messageSent.Username]
		messageToBeSent := types.BroadCast{
			TypeOfMessage: types.Disconnect,
			Message:       messageSent.Username,
			Sender:        messageSent.Username,
			Color:         messageSent.Color,
		}
		bytesSendingMessage, _ := json.Marshal(messageToBeSent)
		if connectionToBeClosed == connection {
			delete(wsManager.RoomWithPeople[messageSent.Room], messageSent.Username)
			connection.Close()
			if len(wsManager.RoomWithPeople[messageSent.Room]) == 0 {
				delete(wsManager.RoomWithPeople, messageSent.Room)
			} else {
				for _, connection := range roomToBeBroadcastedIn {
					connection.WriteMessage(messageType, bytesSendingMessage)
				}
			}
		}
	} else {
		_ = connection.Close()
	}
	wsManager.Mutex.Unlock()
}

func (wsManager *WebSocketManager) SendTextMessage(messageSent types.Message, connection *websocket.Conn, messageType int) {
	wsManager.Mutex.RLock()
	defer wsManager.Mutex.RUnlock()
	roomToBroadCastIn, roomExists := wsManager.RoomWithPeople[messageSent.Room]
	if roomExists {
		var sentMessage types.TextMessageSent
		err := json.Unmarshal(messageSent.Message, &sentMessage)
		if err != nil {
			return
		}
		// check if the user if part of the room
		_, userExists := roomToBroadCastIn[messageSent.Username]
		textMessageSending := types.BroadCast{
			TypeOfMessage: types.TextMessage,
			Message:       sentMessage.Message,
			Sender:        messageSent.Username,
			Color:         messageSent.Color,
		}
		bytesSendingMessage, _ := json.Marshal(textMessageSending)
		if userExists {
			for username, connection := range roomToBroadCastIn {
				if username != messageSent.Username {
					connection.WriteMessage(messageType, bytesSendingMessage)
				}
			}
		}
	}
}

func (wsManager *WebSocketManager) SendPositionMessage(messageSent types.Message, connection *websocket.Conn, messageType int) {
	var positionsToBeBroadcasted types.PositionMessageSent
	err := json.Unmarshal(messageSent.Message, &positionsToBeBroadcasted)
	if err != nil {
		return
	}
	positionsToBeBroadcastedString, err := json.Marshal(positionsToBeBroadcasted)
	if err != nil {
		return
	}
	// TODO:: Here, a check should be applied if the positions provided is correct or not, i.e., accessible or not if not early return
	wsManager.Mutex.RLock()
	defer wsManager.Mutex.RUnlock()
	roomToBroadCastIn, roomExists := wsManager.RoomWithPeople[messageSent.Room]
	if roomExists {
		_, userExists := roomToBroadCastIn[messageSent.Username]
		if userExists {
			positionsMessageSending := types.BroadCast{
				TypeOfMessage: types.PositionMessage,
				Message:       string(positionsToBeBroadcastedString),
				Sender:        messageSent.Username,
				Color:         messageSent.Color,
			}
			bytesSendingMessage, _ := json.Marshal(positionsMessageSending)
			for username, connection := range roomToBroadCastIn {
				if username != messageSent.Username {
					connection.WriteMessage(messageType, bytesSendingMessage)
				}
			}
		}
	}
}

func (wsManager *WebSocketManager) TypeChecker(messageSentInBytes []byte) (bool, string) {
	var messageSent types.Message
	err := json.Unmarshal(messageSentInBytes, &messageSent)
	if err != nil || messageSent.TypeOfMessage == "" || messageSent.Room == "" || messageSent.Username == "" || messageSent.Color == "" {
		return false, ""
	}
	switch messageSent.TypeOfMessage {
	case types.Conect:
		var connectMessage types.ConectMessageSent
		err := json.Unmarshal(messageSent.Message, &connectMessage)
		if err != nil {
			return false, ""
		}
		if connectMessage.Token == "" {
			return false, ""
		}
		return true, string(types.Conect)
	case types.Disconnect:
		return true, string(types.Disconnect)
	case types.TextMessage:
		var textMessage types.TextMessageSent
		err := json.Unmarshal(messageSent.Message, &textMessage)
		if err != nil {
			return false, ""
		}
		if textMessage.Message == "" {
			return false, ""
		}
		return true, string(types.TextMessage)
	case types.PositionMessage:
		var positionMessage types.PositionMessageSent
		err := json.Unmarshal(messageSent.Message, &positionMessage)
		if err != nil {
			return false, ""
		}
		if positionMessage.X == "" || positionMessage.Y == "" {
			return false, ""
		}
		return true, string(types.PositionMessage)
	case types.InitiateCallRequest:
		var initiateCallRequestMessage types.InitiateCallRequestMessage
		err := json.Unmarshal(messageSent.Message, &initiateCallRequestMessage)
		if err != nil {
			return false, ""
		}
		if initiateCallRequestMessage.Receiver == "" {
			return false, ""
		}
		return true, string(types.InitiateCallRequest)
	case types.AcceptCallResponse:
		var acceptCallResponseMessage types.AcceptCallFromReceiver
		err := json.Unmarshal(messageSent.Message, &acceptCallResponseMessage)
		if err != nil {
			return false, ""
		}
		if acceptCallResponseMessage.Initiator == "" {
			return false, ""
		}
		return true, string(types.AcceptCallResponse)
	default:
		return false, ""
	}
}

func (wsManager *WebSocketManager) CleanUp(conn *websocket.Conn) {
	wsManager.Mutex.Lock()
	defer wsManager.Mutex.Unlock()
	for roomId, roomsMaps := range wsManager.RoomWithPeople {
		for username, websocketConn := range roomsMaps {
			if websocketConn == conn {
				delete(wsManager.RoomWithPeople[roomId], username)
				conn.Close()
				if len(wsManager.RoomWithPeople[roomId]) == 0 {
					delete(wsManager.RoomWithPeople, roomId)
				}
				return
			}
		}
	}
}

func (wsManager *WebSocketManager) HandleInitiateCallMessage(messageSent types.Message, connection *websocket.Conn, messageType int) {
	var initiateCallMessage types.InitiateCallRequestMessage
	err := json.Unmarshal(messageSent.Message, &initiateCallMessage)
	if err != nil {
		return
	}
	wsManager.Mutex.RLock()
	defer wsManager.Mutex.RUnlock()
	roomBeingUsed, roomExists := wsManager.RoomWithPeople[messageSent.Room]
	if roomExists {
		_, senderExists := roomBeingUsed[messageSent.Username]
		receiver, receiverExists := roomBeingUsed[initiateCallMessage.Receiver]
		if senderExists && receiverExists {
			messageToBeSent := types.BroadCast{
				TypeOfMessage: types.InitiateCallRequest,
				Sender:        messageSent.Username,
				Color:         messageSent.Color,
				Message:       "",
			}
			messageInBytes, err := json.Marshal(messageToBeSent)
			if err == nil {
				receiver.WriteMessage(messageType, messageInBytes)
			}
		}
	}
}

func (wsManager *WebSocketManager) HandleAcceptCallMessage(messageSent types.Message, connection *websocket.Conn, messageType int) {
	var acceptCallMessage types.AcceptCallFromReceiver
	err := json.Unmarshal(messageSent.Message, &acceptCallMessage)
	if err != nil {
		return
	}
	wsManager.Mutex.RLock()
	defer wsManager.Mutex.RUnlock()
	roomBeingUsed, roomExists := wsManager.RoomWithPeople[messageSent.Room]
	if roomExists {
		initiator, initiatorExists := roomBeingUsed[acceptCallMessage.Initiator]
		_, acceptorExists := roomBeingUsed[messageSent.Username]
		if initiatorExists && acceptorExists {
			messageToBeSent := types.BroadCast{
				TypeOfMessage: types.AcceptCallResponse,
				Sender:        messageSent.Username,
				Color:         messageSent.Color,
				Message:       "",
			}
			messageInBytes, err := json.Marshal(messageToBeSent)
			if err == nil {
				initiator.WriteMessage(messageType, messageInBytes)
			}
		}
	}
}
