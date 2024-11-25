package types

import "encoding/json"

type VideoType string

const (
	IceCandidateMessage    VideoType = "IceCandidates"
	AddPersonToRoomMessage VideoType = "AddToRoom"
	CreateRoomMessage      VideoType = "CreateRoom"
	CreateOfferMessage     VideoType = "CreateOffer"
	CreateAnswerMessage    VideoType = "CreateAnswer"
)

type VideoMessage struct {
	Username      string
	Room          string
	TypeOfMessage VideoType
	Message       json.RawMessage
}

type IceCandidate struct {
	IceCandidate string
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
