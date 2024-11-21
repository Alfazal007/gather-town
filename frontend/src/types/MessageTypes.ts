export enum MessageType {
    PositionMessage = "Position",
    TextMessage = "Text",
    Disconnect = "Disconnect",
    Connect = "Connect"
}

export type Message = {
    Username: string
    Room: string
    TypeOfMessage: MessageType
    Message: TextMessageSent | PositionMessageSent | ConectMessageSent
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
