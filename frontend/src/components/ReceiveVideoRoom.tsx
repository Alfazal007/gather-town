import { useContext, useEffect, useRef, useState } from "react"
import { Card } from "./ui/card"
import { Button } from "./ui/button"
import { Mic, MicOff, Phone, Video, VideoOff } from "lucide-react"
import { UserContext } from "@/context/UserContext"
import { useNavigate, useParams } from "react-router-dom"
import { BroadCastVideoInfo, BroadCastVideoType, VideoMessage, VideoType, CreateRoom, JoinRoom, SDPType, Sdp, IceCandidate } from "@/types/VideoTypes"
import ConnectingToCall from "./ConnectingCall"

const ReceiveVideoRoom = () => {
    const [isMuted, setIsMuted] = useState(false)
    const [isVideoOn, setIsVideoOn] = useState(true)
    const [startedCall, setStartedCall] = useState(false)
    const { user, setUser } = useContext(UserContext)
    const navigate = useNavigate()
    const VIDEOCALL_SOCKET = "ws://localhost:8001/video"
    const { sender, receiver } = useParams()
    const [socket, setSocket] = useState<WebSocket | null>(null)
    const startedCallRef = useRef(startedCall);
    const pcConn = new RTCPeerConnection()

    useEffect(() => {
        startedCallRef.current = startedCall;
    }, [startedCall])

    useEffect(() => {
        if (!user || !sender || !receiver) {
            navigate("/")
            return
        }

        if (user.username != receiver) {
            navigate("/")
            return
        }
        let timeout1: NodeJS.Timeout;
        let timeout2: NodeJS.Timeout;
        const ws = new WebSocket(VIDEOCALL_SOCKET);
        ws.onopen = () => {
            setSocket(ws);
            const joinRoomInternalMessage: JoinRoom = {
                Room: sender + receiver,
                Sender: sender,
                Token: user.accessToken
            }
            const joinRoom: VideoMessage = {
                Room: sender + receiver,
                Username: user.username,
                TypeOfMessage: VideoType.JoinRoomMessage,
                Message: joinRoomInternalMessage
            }
            timeout1 = setTimeout(() => {
                ws.send(JSON.stringify(joinRoom))
            }, 5000)
            timeout2 = setTimeout(() => {
                if (!startedCallRef.current) {
                    console.log({ startedCall })
                    console.log("going back")
                    socket?.close()
                    navigate("/")
                }
            }, 20000)
        };
        ws.onmessage = async (event) => {
            const message: VideoMessage = JSON.parse(event.data)
            if (message.TypeOfMessage == VideoType.SDPRoomMessage) {
                const messageData = message.Message as Sdp
                await pcConn.setRemoteDescription(messageData.Data)
                const answer = await pcConn.createAnswer()
                await pcConn.setLocalDescription(answer)
                const answerDataInternal: Sdp = {
                    Message: SDPType.CreateAnswer,
                    Data: answer
                }
                const answerToSend: VideoMessage = {
                    Message: answerDataInternal,
                    TypeOfMessage: VideoType.SDPRoomMessage,
                    Room: sender + receiver,
                    Username: user.username
                }
                ws.send(JSON.stringify(answerToSend))
            } else if (message.TypeOfMessage == VideoType.IceCandidateMessage) {
                const data: IceCandidate = message.Message as IceCandidate
                await pcConn.addIceCandidate(data.IceCandidate)
            }
        }

        ws.onclose = () => {
            ws.send(JSON.stringify(disconnectMessage))
            setSocket(null);
            navigate("/")
        };

        const disconnectMessage: VideoMessage = {
            TypeOfMessage: VideoType.DisconnectMessage,
            Room: sender + receiver,
            Username: user.username,
            Message: {}
        }
        startReceiving()
        return () => {
            clearTimeout(timeout1);
            clearTimeout(timeout2);
            ws.send(JSON.stringify(disconnectMessage))
        };
    }, [])
    function startReceiving() {
        const video1 = document.createElement('video');
        video1.autoplay = true;
        document.body.appendChild(video1);

        pcConn.ontrack = async (event) => {
            video1.srcObject = new MediaStream([event.track]);
            await video1.play();
        };
    }
    useEffect(() => {
        if (!socket) {
            return
        }

        socket.onmessage = (event) => {
            const message: BroadCastVideoInfo = JSON.parse(event.data);
            switch (message.TypeOfMessage) {
                case BroadCastVideoType.JoinRoomResponse:
                    setStartedCall(true)
                    break
                case BroadCastVideoType.IceCandidates:
                    break
                case BroadCastVideoType.SDP:
                    break
            }
        }

    }, [socket])

    return (
        startedCall ?
            < div className="flex flex-col h-screen bg-gray-100" >
                <main className="flex-1 p-4 relative">
                    {/* Main participant video */}
                    <div className="h-full">
                        <div className="absolute inset-0 bg-black rounded-lg overflow-hidden">
                            <video className="w-full h-full object-cover" autoPlay muted loop>
                                <source src="/placeholder.svg?height=720&width=1280" type="video/mp4" />
                                Your browser does not support the video tag.
                            </video>
                        </div>
                        <div className="absolute bottom-4 left-4 text-white bg-black bg-opacity-50 px-2 py-1 rounded">
                            {sender}
                        </div>
                    </div>

                    <Card className="absolute top-4 right-4 w-48 h-36 overflow-hidden">
                        <video className="w-full h-full object-cover" autoPlay muted loop>
                            Your browser does not support the video tag.
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
                        onClick={() => setIsMuted(!isMuted)}
                        className="bg-white hover:bg-gray-100"
                    >
                        {isMuted ? <MicOff className="h-4 w-4" /> : <Mic className="h-4 w-4" />}
                    </Button>
                    <Button
                        variant="outline"
                        size="icon"
                        onClick={() => setIsVideoOn(!isVideoOn)}
                        className="bg-white hover:bg-gray-100"
                    >
                        {isVideoOn ? <Video className="h-4 w-4" /> : <VideoOff className="h-4 w-4" />}
                    </Button>
                    <Button variant="destructive" size="icon">
                        <Phone className="h-4 w-4" />
                    </Button>
                </div>
            </div > : <div>
                <ConnectingToCall />
            </div>
    )

}

export default ReceiveVideoRoom
