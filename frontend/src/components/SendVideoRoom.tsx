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
    const VIDEOCALL_SOCKET = "ws://192.168.194.11:8001/video"
    //    const VIDEOCALL_SOCKET = "ws://localhost:8001/video"
    const { sender, receiver } = useParams()
    const [socket, setSocket] = useState<WebSocket | null>(null)
    const startedCallRef = useRef(startedCall);
    const pc = new RTCPeerConnection()
    const sendVideoRef = useRef<HTMLVideoElement | null>(null);
    const receiveVideoRef = useRef<HTMLVideoElement | null>(null);
    const audioRef = useRef<HTMLAudioElement | null>(null);
    useEffect(() => {
        startedCallRef.current = startedCall;
    }, [startedCall])

    useEffect(() => {
        return () => {
            pc.close()
        }
    }, [])

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

    async function initializePC() {
        pc.onnegotiationneeded = async () => {
            const offer = await pc.createOffer()
            await pc.setLocalDescription(offer)
            const offerMessageInternal: Sdp = {
                Message: SDPType.CreateOffer,
                Data: offer
            }
            const offerMessageExternal: VideoMessage = {
                Message: offerMessageInternal,
                Room: sender as string + receiver as string,
                TypeOfMessage: VideoType.SDPRoomMessage,
                Username: user?.username as string
            }
            socket?.send(JSON.stringify(offerMessageExternal))
        }
        pc.onicecandidate = (event) => {
            const candidate = event.candidate;
            if (candidate) {
                const iceCandidateInternalMessage: IceCandidate = {
                    IceCandidate: candidate
                }
                const iceCandidateExternalMessage: VideoMessage = {
                    Username: user?.username as string,
                    TypeOfMessage: VideoType.IceCandidateMessage,
                    Room: sender as string + receiver as string,
                    Message: iceCandidateInternalMessage
                }
                socket?.send(JSON.stringify(iceCandidateExternalMessage))
            }
        }
        getCameraStreamAndSend()
    }

    const getCameraStreamAndSend = () => {
        if (navigator.mediaDevices) {
            // Request both video and audio in a single getUserMedia call
            navigator.mediaDevices.getUserMedia({
                video: true,
                audio: true
            }).then((stream) => {
                // Assign the video stream to the sendVideoRef
                if (sendVideoRef.current) {
                    sendVideoRef.current.srcObject = stream;
                    sendVideoRef.current.muted = true;
                    sendVideoRef.current.play().catch((err) => { console.log("Error playing video:", err); });
                }

                // Add both video and audio tracks to the peer connection
                stream.getTracks().forEach((track) => {
                    pc?.addTrack(track);
                });
            }).catch((err) => {
                console.log("Error accessing media devices:", err);
            });
        }

        pc.ontrack = async (event) => {
            if (event.track.kind === 'video') {
                // Handle receiving video
                if (receiveVideoRef.current) {
                    const mediaStream = new MediaStream([event.track]);
                    receiveVideoRef.current.srcObject = mediaStream;
                    receiveVideoRef.current.muted = true;
                    try {
                        receiveVideoRef.current.play().catch((err) => { console.log("Error playing received video:", err); });
                    } catch (err) {
                        console.log("Error in video play:", err);
                    }
                }
            } else if (event.track.kind === 'audio') {
                // Handle receiving audio
                if (audioRef.current) {
                    const mediaStream = new MediaStream([event.track]);
                    audioRef.current.srcObject = mediaStream;
                    try {
                        audioRef.current.play().catch((err) => { console.log("Error playing received audio:", err); });
                    } catch (err) {
                        console.log("Error in audio play:", err);
                    }
                }
            }
        };
    }

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
                    await initializePC()
                    break
                case BroadCastVideoType.SDP:
                    const sdpMessage = message.Message as Sdp
                    if (sdpMessage.Message == SDPType.CreateAnswer) {
                        await pc.setRemoteDescription(sdpMessage.Data)
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
                            <video ref={receiveVideoRef} className="w-full h-full object-cover" />
                        </div>
                        <div className="absolute bottom-4 left-4 text-white bg-black bg-opacity-50 px-2 py-1 rounded">
                            {receiver}
                        </div>
                    </div>

                    <Card className="absolute top-4 right-4 w-48 h-36 overflow-hidden">
                        <video ref={sendVideoRef} className="w-full h-full object-cover">
                        </video>
                        <div className="absolute bottom-2 left-2 text-white bg-black bg-opacity-50 px-2 py-1 rounded text-sm">
                            You
                        </div>
                    </Card>
                </main>
                <audio ref={audioRef}></audio>
                {/* Call controls */}
                <div className="p-4 flex justify-center space-x-4">
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
