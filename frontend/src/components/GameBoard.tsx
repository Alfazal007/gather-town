import { useContext, useEffect, useState } from "react";
import { Player } from "./Player";
import { UserContext } from "@/context/UserContext";
import { useNavigate, useParams } from "react-router-dom";
import { BroadCast, Message, MessageType, PositionMessageSent } from "@/types/MessageTypes";
import { Button } from "./ui/button";
import ChatDisplay from "./ChatDisplay";
import VideoCall from "./VideoCallButtons";
import { CustomModal } from "./CustomModal";

export type OtherPlayersType = {
    [key: string]: {
        username: string,
        color: string,
        x: number,
        y: number
    };
};

type TextMessageType = {
    sender: string,
    message: string
}

const colors = [
    "bg-gray-600", "bg-gray-700",
    "bg-blue-600", "bg-blue-700",
    "bg-indigo-600", "bg-indigo-700",
    "bg-cyan-600", "bg-cyan-700",
    "bg-red-600", "bg-red-700",
    "bg-yellow-600", "bg-yellow-700",
    "bg-amber-600", "bg-amber-700",
    "bg-green-600", "bg-green-700",
    "bg-teal-600", "bg-teal-700",
    "bg-pink-600", "bg-pink-700"
];

const GameBoard = () => {
    const [otherPlayers, setOtherPlayers] = useState<OtherPlayersType>({});
    const [position, setPosition] = useState({ x: 0, y: 0 });
    const [_, setCallIncoming] = useState(false);
    const [callComingFrom, setCallComingFrom] = useState("")
    const boardWidth = 1400;
    const boardHeight = 900;
    const characterSize = 40;
    const characterHeight = 100;
    const { user } = useContext(UserContext)
    const navigate = useNavigate()
    const { roomId } = useParams()
    //    const WS_URL = "ws://localhost:8001/ws";
    const WS_URL = "ws://192.168.194.11:8001/ws";
    const [color, setColor] = useState("bg-red-600")
    const [messages, setMessages] = useState<TextMessageType[]>([])
    const [moveStarted, setMoveStarted] = useState<boolean>(false)
    const disconnectMessage: Message = {
        Color: color,
        Username: user?.username as string,
        TypeOfMessage: MessageType.Disconnect,
        Room: roomId as string,
        Message: {
            Message: "Disconnect"
        }
    }
    const [acceptor, setAcceptor] = useState<string>("")

    const [isModalOpen, setIsModalOpen] = useState(false)
    const handleOpenModal = () => setIsModalOpen(true)
    const handleCloseModal = () => setIsModalOpen(false)


    function generateRandomColor() {
        const index = Math.floor(Math.random() * (colors.length));
        setColor(colors[index])
    }

    const [socket, setSocket] = useState<WebSocket | null>(null);

    useEffect(() => {
        const ws = new WebSocket(WS_URL);
        ws.onopen = () => {
            setSocket(ws);
        };
        ws.onclose = () => {
            setTimeout(() => {
                setSocket(null);
                //navigate("/")
            }, 100);

        };
        return () => {
        };
    }, []);

    const handleAccept = () => {
        acceptCall()
    }


    const updatePositionOfOtherPlayer = (x: string, y: string, username: string, color: string) => {
        const numberX = parseInt(x) || 0
        const numberY = parseInt(y) || 0
        setOtherPlayers(prev => ({
            ...prev,
            [username]: {
                username,
                color,
                x: numberX,
                y: numberY
            }
        }))
    }

    const removeOtherPlayer = (username: string) => {
        setOtherPlayers(prev => {
            const updatedHashmap = { ...prev }
            delete updatedHashmap[username]
            return updatedHashmap
        })
    }
    useEffect(() => {
        generateRandomColor()
    }, [])

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            setPosition(prev => {
                let newX = prev.x;
                let newY = prev.y;

                if (moveStarted) {
                    switch (e.key) {
                        case 'ArrowUp':
                            newY = Math.max(0, prev.y - 10);
                            break;
                        case 'ArrowDown':
                            newY = Math.min(boardHeight - characterSize, prev.y + 10);
                            break;
                        case 'ArrowLeft':
                            newX = Math.max(0, prev.x - 10);
                            break;
                        case 'ArrowRight':
                            newX = Math.min(boardWidth - characterHeight, prev.x + 10);
                            break;
                    }
                }
                return { x: newX, y: newY };
            });
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, [moveStarted]);

    async function initialize() {
        if (!socket || !roomId || !user) {
            socket?.send(JSON.stringify(disconnectMessage))
            socket?.close()
            return
        }
        const connectMessage: Message = {
            Room: roomId as string,
            Color: color,
            TypeOfMessage: MessageType.Connect,
            Username: user?.username as string,
            Message: { Token: user?.accessToken as string }
        }
        setMoveStarted(true)
        socket?.send(JSON.stringify(connectMessage))
    }

    useEffect(() => {
        if (!socket || !roomId || !user) {
            return
        }
        let positionData: Message = {
            Color: color,
            TypeOfMessage: MessageType.PositionMessage,
            Room: roomId,
            Username: user?.username,
            Message: {
                X: position.x + "",
                Y: position.y + ""
            }
        }
        socket.send(JSON.stringify(positionData))
    }, [position])

    function sendTextMessage(message: string) {
        if (!socket || !roomId || !user) {
            return
        }

        let messageData: Message = {
            Color: color,
            TypeOfMessage: MessageType.TextMessage,
            Room: roomId,
            Username: user?.username,
            Message: {
                Message: message
            }
        }
        socket.send(JSON.stringify(messageData))
        setMessages((prev) => {
            return [...prev, { sender: user.username, message: message }]
        })
    }

    function acceptCall() {
        handleCloseModal()
        if (!socket || !roomId || !user) {
            return
        }

        let messageData: Message = {
            Color: color,
            TypeOfMessage: MessageType.AcceptCallResponse,
            Room: roomId,
            Username: user?.username,
            Message: {
                Initiator: callComingFrom
            }
        }
        socket.send(JSON.stringify(messageData))

        socket.send(JSON.stringify(disconnectMessage));
        navigate(`/video/${callComingFrom}/${user.username}/receiver`)
    }

    function disconnectAndMoveToCallAttendScreen(username: string) {
        if (!socket || !roomId || !user) {
            return
        }
        socket.send(JSON.stringify(disconnectMessage));
        navigate(`/video/${user.username}/${username}/sender`)
    }

    function someOneRequestedACallHandler(requestor: string) {
        setCallComingFrom(requestor);
        setCallIncoming(true);
        handleOpenModal()
    }

    function makeCall(receiver: string) {
        if (!user || !roomId || !socket) {
            return
        }
        let messageToBeSent: Message = {
            Room: roomId,
            Username: user.username,
            Color: color,
            TypeOfMessage: MessageType.InitiateCallRequest,
            Message: {
                Sender: user.username,
                Receiver: receiver
            }
        }
        socket.send(JSON.stringify(messageToBeSent))
    }

    function someoneAcceptedACallHandler(acceptor: string) {
        setAcceptor(acceptor)
    }

    useEffect(() => {
        if (!user) {
            socket?.send(JSON.stringify(disconnectMessage))
            socket?.close()
            //navigate("/")
            return
        }
        if (!socket) {
            return
        }
        socket.onmessage = (event) => {
            const message: BroadCast = JSON.parse(event.data);
            switch (message.TypeOfMessage) {
                case MessageType.PositionMessage:
                    let positionsOfOtherUser: PositionMessageSent = JSON.parse(message.Message)
                    updatePositionOfOtherPlayer(positionsOfOtherUser.X + "", positionsOfOtherUser.Y + "", message.Sender, message.Color)
                    break
                case MessageType.TextMessage:
                    setMessages((prev) => {
                        return [...prev, { sender: message.Sender, message: message.Message }]
                    })
                    break
                case MessageType.Disconnect:
                    removeOtherPlayer(message.Sender)
                    break
                case MessageType.InitiateCallRequest:
                    someOneRequestedACallHandler(message.Sender)
                    break
                case MessageType.AcceptCallResponse:
                    someoneAcceptedACallHandler(message.Sender)
                    break
            }
        }
    }, [socket])

    return (
        <>
            <div className="flex mt-2">
                <div
                    className="relative bg-customBackgroundGame ml-3"
                    style={{
                        width: `${boardWidth}px`,
                        height: `${boardHeight}px`
                    }}
                >
                    <CustomModal
                        isOpen={isModalOpen}
                        onClose={handleCloseModal}
                        onAccept={handleAccept}
                        caller={callComingFrom}
                    />
                    <Player x={position.x} y={position.y} color={color} username={user?.username as string} />
                    {
                        Object.entries(otherPlayers).map(([id, player]) => (
                            <Player key={id} y={player.y} x={player.x} color={player.color} username={player.username} />
                        ))
                    }

                </div>
                {
                    moveStarted &&
                    <div>
                        <ChatDisplay otherPeople={otherPlayers} messages={messages}
                            onLeaveRoom={() => { socket?.send(JSON.stringify(disconnectMessage)); socket?.close(); setSocket(null); navigate("/") }} onSendMessage={(message: string) => { sendTextMessage(message) }}
                        />
                        <VideoCall otherPeople={otherPlayers} makeACall={makeCall} acceptor={acceptor} disconnectAndMoveToCallAttendScreen={disconnectAndMoveToCallAttendScreen} />
                    </div>
                }
                <div className="flex items-center">
                    <div>
                        {
                            !moveStarted &&
                            <Button onClick={initialize} className="ml-4 cursor-pointer">Join Room</Button>
                        }

                        {
                            !moveStarted &&
                            <Button onClick={() => { socket?.close(); setSocket(null); navigate("/") }} className="ml-4 cursor-pointer">Go back</Button>
                        }
                    </div>
                </div>
            </div>
        </>
    );
};

export default GameBoard;
