import { useContext, useEffect, useRef, useState } from "react"
import { Card } from "./ui/card"
import { Button } from "./ui/button"
import { Mic, MicOff, Phone, Video, VideoOff } from "lucide-react"
import { UserContext } from "@/context/UserContext"
import { useNavigate, useParams } from "react-router-dom"
import { BroadCastVideoInfo, BroadCastVideoType, VideoMessage, VideoType, CreateRoom, IceCandidate, Sdp, SDPType } from "@/types/VideoTypes"
import ConnectingToCall from "./ConnectingCall"

const SendVideoRoom = () => {
    const [isMuted, setIsMuted] = useState(false)
    const [isVideoOn, setIsVideoOn] = useState(false)
    const [startedCall, setStartedCall] = useState(false)
    const { user } = useContext(UserContext)
    const navigate = useNavigate()
    const VIDEOCALL_SOCKET = "ws://localhost:8001/video"
    const { sender, receiver } = useParams()
    const [socket, setSocket] = useState<WebSocket | null>(null)
    const startedCallRef = useRef(startedCall);

    useEffect(() => {
        startedCallRef.current = startedCall;
    }, [startedCall])

    useEffect(() => {
        if (!user || !sender || !receiver) {
            navigate("/")
            return
        }

        if (user.username != sender) {
            navigate("/")
            return
        }
        let timeout: NodeJS.Timeout

        const ws = new WebSocket(VIDEOCALL_SOCKET);

        ws.onopen = () => {
            setSocket(ws);
            const createRoomData: CreateRoom = {
                Token: user.accessToken,
                Sender: user.username,
                Receiver: receiver
            }
            const connectToRoom: VideoMessage = {
                TypeOfMessage: VideoType.CreateRoomMessage,
                Room: sender + receiver,
                Username: sender,
                Message: createRoomData
            }

            ws.send(JSON.stringify(connectToRoom))
            timeout = setTimeout(() => {
                if (!startedCallRef.current) {
                    ws.close()
                    navigate("/")
                }
            }, 20000)
        };

        const disconnectMessage: VideoMessage = {
            TypeOfMessage: VideoType.DisconnectMessage,
            Room: sender + receiver,
            Username: user.username,
            Message: {}
        }
        setSocket(ws)
        ws.onclose = () => {
            ws.send(JSON.stringify(disconnectMessage))
            setSocket(null);
            navigate("/")
        };

        return () => {
            ws.send(JSON.stringify(disconnectMessage))
            clearTimeout(timeout)
        }
    }, [])

    useEffect(() => {
        if (!socket) {
            return
        }

        socket.onmessage = async (event) => {
            const message: BroadCastVideoInfo = JSON.parse(event.data);
            switch (message.TypeOfMessage) {
                case BroadCastVideoType.CreateRoomResponse:
                    break
                case BroadCastVideoType.JoinRoomResponse:
                    setStartedCall(true)
                    break
                case BroadCastVideoType.IceCandidates:
                    const iceCandidates: IceCandidate = message.Message as IceCandidate
                    break
                case BroadCastVideoType.SDP:
                    const sdpMessage = message.Message as Sdp
                    if (sdpMessage.Message == SDPType.CreateAnswer) {
                    } else if (sdpMessage.Message == SDPType.CreateOffer) {
                    }
                    break
            }
        }
    }, [socket])

    return (
        startedCall ?
            <div className="flex flex-col h-screen bg-gray-100">
                <main className="flex-1 p-4 relative">
                    {/* Main participant video */}
                    <div className="h-full">
                        <div className="absolute inset-0 bg-black rounded-lg overflow-hidden">
                            <video className="w-full h-full object-cover" autoPlay muted loop>
                            </video>
                        </div>
                        <div className="absolute bottom-4 left-4 text-white bg-black bg-opacity-50 px-2 py-1 rounded">
                            {receiver}
                        </div>
                    </div>

                    <Card className="absolute top-4 right-4 w-48 h-36 overflow-hidden">
                        <video className="w-full h-full object-cover" autoPlay muted loop>
                        </video>
                        <div className="absolute bottom-2 left-2 text-white bg-black bg-opacity-50 px-2 py-1 rounded text-sm">
                            You
                        </div>
                    </Card>
                </main>

                {/* Call controls */}
                <div className="p-4 flex justify-center space-x-4">
                    <Button
                        variant="outline"
                        size="icon"
                        className="bg-white hover:bg-gray-100"
                    >
                        {isMuted ? <MicOff className="h-4 w-4" /> : <Mic className="h-4 w-4" />}
                    </Button>
                    <Button
                        variant="outline"
                        size="icon"
                        className="bg-white hover:bg-gray-100"
                    >
                        {isVideoOn ? <Video className="h-4 w-4" /> : <VideoOff className="h-4 w-4" />}
                    </Button>
                    <Button variant="destructive" size="icon">
                        <Phone className="h-4 w-4" />
                    </Button>
                </div>
            </div> : <div>
                <ConnectingToCall />
            </div>
    )

}

export default SendVideoRoom
