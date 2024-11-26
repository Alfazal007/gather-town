package managers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Alfazal007/gather-town/types"
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

func (vsManager *VideoRoomManager) CleanUp(conn *websocket.Conn) {
	vsManager.Mutex.Lock()
	defer vsManager.Mutex.Unlock()
	for roomId, roomsMaps := range vsManager.RoomWithTwoPeople {
		if roomsMaps.Person1.Conn == conn || roomsMaps.Person2.Conn == conn || roomsMaps.Person1.Conn == nil || roomsMaps.Person2.Conn == nil {
			if roomsMaps.Person2.Conn != nil {
				_ = roomsMaps.Person2.Conn.Close()
			}
			if roomsMaps.Person1.Conn != nil {
				_ = roomsMaps.Person1.Conn.Close()
			}
			delete(vsManager.RoomWithTwoPeople, roomId)
			return
		}
	}
}

func (vsManager *VideoRoomManager) TypeChecker(messageSentInBytes []byte) (bool, string) {
	var videoMessageSent types.VideoMessage
	err := json.Unmarshal(messageSentInBytes, &videoMessageSent)
	if err != nil || videoMessageSent.TypeOfMessage == "" || videoMessageSent.Room == "" || videoMessageSent.Username == "" {
		return false, ""
	}
	switch videoMessageSent.TypeOfMessage {
	case types.CreateRoomMessage:
		var createRoomMessage types.CreateRoom
		err := json.Unmarshal(videoMessageSent.Message, &createRoomMessage)
		if err != nil {
			return false, ""
		}
		if createRoomMessage.Sender == "" || createRoomMessage.Receiver == "" || createRoomMessage.Token == "" {
			return false, ""
		}
		return true, string(types.CreateRoomMessage)
	case types.IceCandidateMessage:
		var iceCandidateMessage types.IceCandidate
		err := json.Unmarshal(videoMessageSent.Message, &iceCandidateMessage)
		if err != nil {
			return false, ""
		}
		return true, string(types.IceCandidateMessage)
	case types.SDPRoomMessage:
		var sdpTypeMessage types.Sdp
		err := json.Unmarshal(videoMessageSent.Message, &sdpTypeMessage)
		if err != nil {
			return false, ""
		}
		if sdpTypeMessage.Message != types.CreateAnswer && sdpTypeMessage.Message != types.CreateOffer {
			return false, ""
		}
		return true, string(types.SDPRoomMessage)
	case types.JoinRoomMessage:
		var joinRoomType types.JoinRoom
		err := json.Unmarshal(videoMessageSent.Message, &joinRoomType)
		if err != nil {
			return false, ""
		}
		if joinRoomType.Room == "" || joinRoomType.Token == "" || joinRoomType.Sender == "" {
			return false, ""
		}
		return true, string(types.JoinRoomMessage)
	case types.DisconnectMessage:
		return true, string(types.Disconnect)
	default:
		return false, ""
	}
}

func (vsManager *VideoRoomManager) RegisterUserForVideo(token string, username string) bool {
	if len(token) < 5 {
		return false
	}
	url := fmt.Sprintf("%v/%v/username/%v", types.UrlToVideo, token, username)
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

func (vsManager *VideoRoomManager) CreateRoomForVideo(message types.VideoMessage, conn *websocket.Conn, messageType int) bool {
	vsManager.Mutex.Lock()
	defer vsManager.Mutex.Unlock()
	var messageOfConnection types.CreateRoom
	err := json.Unmarshal(message.Message, &messageOfConnection)
	if err != nil {
		return false
	}
	// check if either of the users are part of some other room
	for _, twoPeople := range vsManager.RoomWithTwoPeople {
		person1 := twoPeople.Person1
		person2 := twoPeople.Person2
		if (person1.Username == messageOfConnection.Sender && person1.Conn != nil) || (person2.Username == messageOfConnection.Sender && person2.Conn != nil) || (person1.Username == messageOfConnection.Receiver && person1.Conn != nil) || (person2.Username == messageOfConnection.Receiver && person2.Conn != nil) {
			return false
		}
	}
	newKey := fmt.Sprintf("%v%v", messageOfConnection.Sender, messageOfConnection.Receiver)
	vsManager.RoomWithTwoPeople[newKey] = TwoPeople{
		Person1: Person{
			Username: messageOfConnection.Sender,
			Conn:     conn,
		},
		Person2: Person{
			Username: messageOfConnection.Receiver,
			Conn:     nil,
		},
	}
	createdMessage := types.BroadCastVideoInfo{
		Room:          newKey,
		Username:      messageOfConnection.Sender,
		Message:       nil,
		TypeOfMessage: types.CreateRoomResponse,
	}
	dataInBytes, err := json.Marshal(createdMessage)
	if err != nil {
		return false
	}
	conn.WriteMessage(messageType, dataInBytes)
	return true
}

func (vsManager *VideoRoomManager) ForwardIceCandidates(message types.VideoMessage, conn *websocket.Conn, messageType int) {
	vsManager.Mutex.RLock()
	defer vsManager.Mutex.RUnlock()
	var messageOfIceCandidates types.IceCandidate
	err := json.Unmarshal(message.Message, &messageOfIceCandidates)
	if err != nil {
		return
	}

	roomToBeForwardedTo, roomExists := vsManager.RoomWithTwoPeople[message.Room]
	if roomExists {
		if roomToBeForwardedTo.Person1.Conn == nil || roomToBeForwardedTo.Person2.Conn == nil {
			return
		}
		createdMessage := types.BroadCastVideoInfo{
			Room:          message.Room,
			Username:      message.Username,
			Message:       message.Message,
			TypeOfMessage: types.IceCandidates,
		}
		messageInBytes, err := json.Marshal(createdMessage)
		if err != nil {
			return
		}
		if conn == roomToBeForwardedTo.Person1.Conn {
			roomToBeForwardedTo.Person2.Conn.WriteMessage(messageType, messageInBytes)
		}
		if conn == roomToBeForwardedTo.Person2.Conn {
			roomToBeForwardedTo.Person1.Conn.WriteMessage(messageType, messageInBytes)
		}
	}
}

func (vsManager *VideoRoomManager) ForwardSDPData(message types.VideoMessage, conn *websocket.Conn, messageType int) {
	vsManager.Mutex.RLock()
	defer vsManager.Mutex.RUnlock()
	var messageOfIceCandidates types.Sdp
	err := json.Unmarshal(message.Message, &messageOfIceCandidates)
	if err != nil {
		return
	}

	roomToBeForwardedTo, roomExists := vsManager.RoomWithTwoPeople[message.Room]
	if roomExists {
		if roomToBeForwardedTo.Person1.Conn == nil || roomToBeForwardedTo.Person2.Conn == nil {
			return
		}
		createdMessage := types.BroadCastVideoInfo{
			Room:          message.Room,
			Username:      message.Username,
			Message:       message.Message,
			TypeOfMessage: types.SDP,
		}
		messageInBytes, err := json.Marshal(createdMessage)
		if err != nil {
			return
		}
		if conn == roomToBeForwardedTo.Person1.Conn {
			roomToBeForwardedTo.Person2.Conn.WriteMessage(messageType, messageInBytes)
		} else {
			roomToBeForwardedTo.Person1.Conn.WriteMessage(messageType, messageInBytes)
		}
	}
}

func (vsManager *VideoRoomManager) JoinRoomBySecondPerson(message types.VideoMessage, conn *websocket.Conn, messageType int) {
	vsManager.Mutex.Lock()
	defer vsManager.Mutex.Unlock()
	room, roomExists := vsManager.RoomWithTwoPeople[message.Room]
	if roomExists && room.Person2.Conn == nil && room.Person2.Username == message.Username {
		person1 := room.Person1
		vsManager.RoomWithTwoPeople[message.Room] = TwoPeople{
			Person1: person1,
			Person2: Person{Username: message.Username, Conn: conn},
		}
		createdMessage := types.BroadCastVideoInfo{
			Room:          message.Room,
			Username:      message.Username,
			Message:       message.Message,
			TypeOfMessage: types.JoinRoomResponse,
		}
		messageInBytes, err := json.Marshal(createdMessage)
		if err != nil {
			return
		}
		person1.Conn.WriteMessage(messageType, messageInBytes)
		conn.WriteMessage(messageType, messageInBytes)
	}
}

func (vsManager *VideoRoomManager) DisconnectVideoCall(message types.VideoMessage, conn *websocket.Conn, messageType int) {
	fmt.Println("called")
	vsManager.Mutex.Lock()
	defer vsManager.Mutex.Unlock()
	room, roomExists := vsManager.RoomWithTwoPeople[message.Room]
	if roomExists {
		if room.Person1.Username == message.Username || room.Person2.Username == message.Username {
			if conn == room.Person1.Conn || conn == room.Person2.Conn {
				if room.Person1.Conn != nil {
					_ = room.Person1.Conn.Close()
				}
				if room.Person2.Conn != nil {
					_ = room.Person2.Conn.Close()
				}
				delete(vsManager.RoomWithTwoPeople, message.Room)
			}
		}
	}
}
