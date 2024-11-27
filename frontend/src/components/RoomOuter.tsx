import { UserContext } from "@/context/UserContext"
import axios from "axios"
import { useContext, useEffect, useState } from "react"
import { useNavigate, useParams } from "react-router-dom"
import { Input } from "./ui/input"
import { Button } from "./ui/button"
import Navbar from "./Navbar"
import { Card, CardContent, CardHeader, CardTitle } from "./ui/card"
import { User } from "lucide-react"
import { toast } from "@/hooks/use-toast"

type RoomMemberData = {
    roomMemberId: string,
    roomMemberUsername: string
}

type RoomDataType = {
    roomId: string,
    roomName: string,
    adminId: string,
    roomMembers: RoomMemberData[],
}

type FoundUserDataType = {
    id: string,
    username: string,
    email: string
}

type AdminData = {
    adminName: string,
    adminId: string,
}

const RoomOuter = () => {
    const { user, setUser } = useContext(UserContext)
    const navigate = useNavigate()
    const [isAdmin, setIsAdmin] = useState(false)
    const [adminData, setAdminData] = useState<AdminData | null>(null)
    const [roomData, setRoomData] = useState<RoomDataType | null>()
    const [foundUser, setFoundUser] = useState<boolean>(false)
    const [userToBeAdded, setUserToBeAdded] = useState<string>("");
    const [foundUserData, setFoundUserData] = useState<FoundUserDataType | null>(null)
    const { roomId } = useParams()

    useEffect(() => {
        if (!user) {
            navigate("/signin")
            return
        }
        fetchRoomData()
    }, [user])

    const fetchRoomData = async () => {
        const url = `http://localhost:8000/api/v1/room/roomId/${roomId}`
        const res = await axios.get(
            url,
            {
                headers: {
                    Authorization: `Bearer ${user?.accessToken}`
                }
            }
        );
        if (res.status != 200) {
            navigate("/")
        }
        setRoomData(res.data)
        const adminUrl = `http://localhost:8000/api/v1/user/get-admin/roomId/${roomId}`
        const adminRes = await axios.get(
            adminUrl,
            {
                headers: {
                    Authorization: `Bearer ${user?.accessToken}`
                }
            }
        );
        if (adminRes.status != 200) {
            navigate("/")
        }
        setAdminData(adminRes.data)
        setFoundUser(false)
        setUserToBeAdded("")
        setFoundUserData(null)
    }

    useEffect(() => {
        if (!roomData || !user) {
            return
        }
        if (roomData.adminId == user.id) {
            setIsAdmin(true)
        }
    }, [roomData])

    async function handleAddRoom() {
        const url = "http://localhost:8000/api/v1/room/add-member"
        try {
            const res = await axios.post(
                url,
                {
                    "roomId": roomId,
                    "userId": foundUserData?.id
                },
                {
                    headers: {
                        Authorization: `Bearer ${user?.accessToken}`
                    }
                }
            );
            if (res.status == 200) {
                toast({
                    title: "Added user to this room"
                })
                fetchRoomData()
            } else {
                toast({
                    title: "Error adding the user to this room",
                    variant: "destructive"
                })
                setFoundUser(false)
                setUserToBeAdded("")
                setFoundUserData(null)
            }
        } catch (err) {
            toast({
                title: "Error adding the user to this room",
                variant: "destructive"
            })
            setFoundUser(false)
            setUserToBeAdded("")
            setFoundUserData(null)
        }
        console.log({ foundUserData })
    }

    async function handleFindUser() {
        const url = `http://localhost:8000/api/v1/user/username/${userToBeAdded}`
        try {
            const res = await axios.get(
                url,
                {
                    headers: {
                        Authorization: `Bearer ${user?.accessToken}`
                    }
                }
            );
            if (res.status == 200) {
                setFoundUser(true)
                setUserToBeAdded("")
                setFoundUserData(res.data)

            } else {
                setFoundUser(false)
                setUserToBeAdded("")
                setFoundUserData(null)
            }
        } catch (err) {
            setFoundUser(false)
            setUserToBeAdded("")
            setFoundUserData(null)
        }
    }

    return (
        <>
            <Navbar setUser={setUser} />
            <div className="w-full max-w-sm mx-auto space-y-4">

                {
                    isAdmin &&
                    <div className="flex space-x-2">
                        <Input
                            type="text"
                            placeholder="Enter username of the person to find"
                            value={userToBeAdded}
                            onChange={(e) => { setUserToBeAdded(e.target.value) }}
                            aria-label="Room name or ID"
                        />
                        <Button disabled={!roomId} onClick={handleFindUser}>
                            Search
                        </Button>
                    </div>
                }
                {
                    isAdmin && foundUser && (
                        <Card
                            className="p-4 cursor-pointer"
                            onClick={handleAddRoom}
                        >
                            <div className="flex items-center space-x-4">
                                <User className="h-10 w-10 text-gray-500" />
                                <div>
                                    <h3 className="font-semibold">{foundUserData?.username}</h3>
                                    <p className="text-sm text-gray-500">Click to add this user to this room</p>
                                </div>
                            </div>
                        </Card>
                    )
                }
                <div className="m-2 p-2">
                    <Button className="w-full text-md" onClick={() => navigate(`/room/play/${roomId}`)} disabled={!roomId}>
                        Join Room
                    </Button>
                </div>

            </div >
            <Card className="w-full max-w-md mx-auto">
                <CardHeader>
                    <CardTitle className="text-2xl font-bold text-center">Users in this Room</CardTitle>
                </CardHeader>
                <CardContent>
                    <CardHeader>
                        <CardTitle className="text-xl font-bold text-center">Admin</CardTitle>
                    </CardHeader>
                    <div className="flex items-center space-x-4 p-2 rounded-lg hover:bg-muted transition-colors">
                        <span className="flex-grow">{adminData?.adminName}</span>
                    </div>
                    <CardHeader>
                        <CardTitle className="text-xl font-bold text-center">Members</CardTitle>
                    </CardHeader>
                    <ul className="space-y-4">
                        {roomData?.roomMembers.map((user) => (
                            <li key={user.roomMemberId} className="flex items-center space-x-4 p-2 rounded-lg hover:bg-muted transition-colors">
                                <span className="flex-grow">{user.roomMemberUsername}</span>
                            </li>
                        ))}
                    </ul>
                </CardContent>
            </Card>
        </>
    )
}

export default RoomOuter
