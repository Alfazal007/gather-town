package types

import "encoding/json"

type VideoType string

const (
	IceCandidateMessage VideoType = "IceCandidates"
	CreateRoomMessage   VideoType = "CreateRoom"
	SDPRoomMessage      VideoType = "SDP"
	JoinRoomMessage     VideoType = "JoinRoom"
	DisconnectMessage   VideoType = "Disconnect"
)

type SDPType string

const (
	CreateOffer  SDPType = "CreateOffer"
	CreateAnswer SDPType = "CreateAnswer"
)

type VideoMessage struct {
	Username      string
	Room          string
	TypeOfMessage VideoType
	Message       json.RawMessage
}

type IceCandidate struct {
	IceCandidate json.RawMessage
}

type CreateRoom struct {
	Sender   string
	Receiver string
	Token    string
}

type BroadCastVideoInfo struct {
	Room     string
	Username string
	Message  json.RawMessage
}

type RoomCreationState struct {
	CreatedRoom bool
}

type Sdp struct {
	Message SDPType
	Data    json.RawMessage
}

type JoinRoom struct {
	Sender string
	Room   string
	Token  string
}
