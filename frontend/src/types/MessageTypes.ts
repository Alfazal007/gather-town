export enum MessageType {
    PositionMessage = "Position",
    TextMessage = "Text",
    Disconnect = "Disconnect",
    Connect = "Connect",
    InitiateCallRequest = "IntiateCallRequest",
    AcceptCallResponse = "AcceptCallResponse"
}

export type Message = {
    Username: string
    Room: string
    TypeOfMessage: MessageType
    Message: TextMessageSent | PositionMessageSent | ConectMessageSent | InitiateCallToReceiverFromServer | AcceptCallFromReceiver
    Color: string
}

type TextMessageSent = {
    Message: string
}

export type PositionMessageSent = {
    X: string
    Y: string
}

type ConectMessageSent = {
    Token: string
}

export type BroadCast = {
    TypeOfMessage: MessageType,
    Message: string,
    Sender: string,
    Color: string
}

export type InitiateCallToReceiverFromServer = {
    Sender: string
    Receiver: string
}

export type AcceptCallFromReceiver = {
    Initiator: string
}

