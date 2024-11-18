package types

import "encoding/json"

type MessageType string

const (
	PositionMessage MessageType = "Position"
	TextMessage     MessageType = "Text"
	Disconnect      MessageType = "Disconnect"
	Conect          MessageType = "Connect"
)

type Message struct {
	Username      string
	Room          string
	TypeOfMessage MessageType
	Message       json.RawMessage
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
