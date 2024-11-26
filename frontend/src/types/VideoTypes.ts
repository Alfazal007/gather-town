export enum VideoType {
    IceCandidateMessage = "IceCandidates",
    CreateRoomMessage = "CreateRoom",
    SDPRoomMessage = "SDP",
    JoinRoomMessage = "JoinRoom",
    DisconnectMessage = "Disconnect"
}


export enum SDPType {
    CreateOffer = "CreateOffer",
    CreateAnswer = "CreateAnswer"
}

export type VideoMessage = {
    Username: string,
    Room: string,
    TypeOfMessage: VideoType,
    Message: IceCandidate | CreateRoom | Sdp | JoinRoom | {},
}

export type IceCandidate = {
    IceCandidate: JSON
}

export type CreateRoom = {
    Sender: string,
    Receiver: string,
    Token: string,
}


export enum BroadCastVideoType {
    IceCandidates = "IceCandidates",
    SDP = "SDP",
    CreateRoomResponse = "CreateRoomResponse",
    JoinRoomResponse = "JoinRoomResponse",
}

export type BroadCastVideoInfo = {
    Room: string,
    Username: string,
    Message: JSON,
    TypeOfMessage: BroadCastVideoType,
}

export type RoomCreationState = {
    CreatedRoom: boolean,
}

export type Sdp = {
    Message: SDPType,
    Data: JSON
}

export type JoinRoom = {
    Sender: string,
    Room: string,
    Token: string,
}
