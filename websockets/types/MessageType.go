package types

import "encoding/json"

type MessageType string

const (
	PositionMessage     MessageType = "Position"
	TextMessage         MessageType = "Text"
	Disconnect          MessageType = "Disconnect"
	Conect              MessageType = "Connect"
	InitiateCallRequest MessageType = "IntiateCallRequest"
)

type Message struct {
	Username      string
	Room          string
	TypeOfMessage MessageType
	Message       json.RawMessage
	Color         string
}

type InitiateCallRequestMessage struct {
	Receiver string
}

type TextMessageSent struct {
	Message string
}

type PositionMessageSent struct {
	X string
	Y string
}

type ConectMessageSent struct {
	Token string
}

type BroadCast struct {
	TypeOfMessage MessageType
	Message       string
	Sender        string
	Color         string
}

type InitiateCallToReceiverFromServer struct {
	Sender   string
	Receiver string
}
