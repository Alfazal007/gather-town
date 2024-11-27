import { UserContext } from '@/context/UserContext'
import { useContext, useEffect, useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import Navbar from './Navbar'
import axios from 'axios'
import {
    Card,
    CardContent,
    CardHeader,
    CardTitle,
} from "@/components/ui/card"

type RoomType = {
    id: string,
    roomName: string,
    adminId: string
}

const Landing = () => {
    const navigate = useNavigate()
    const { user, setUser } = useContext(UserContext)
    const [rooms, setRooms] = useState<RoomType[]>([])

    useEffect(() => {
        if (!user) {
            navigate("/signin")
            return
        }
        fetchRoomDetails()
    }, [user])

    const fetchRoomDetails = async () => {
        const url = "http://192.168.194.11:8000/api/v1/user/get-rooms"
        const res = await axios.get(
            url,
            {
                headers: {
                    Authorization: `Bearer ${user?.accessToken}`
                }
            }
        );
        console.log(res.data)
        if (res.status == 200) {
            setRooms(res.data)
        }
    }
    return (
        <>
            <Navbar setUser={setUser} />
            <Card className="w-full max-w-md mx-auto">
                <CardHeader>
                    {rooms.length > 0 && <CardTitle>Your rooms</CardTitle>}
                </CardHeader>
                <CardContent>
                    <ul className="space-y-2">
                        {rooms.map((room, index) => (
                            <Link to={`/room/${room.id}`}>
                                <li
                                    key={index}
                                    className="bg-secondary m-2 text-secondary-foreground rounded-lg p-3 transition-colors hover:bg-secondary/80"
                                >
                                    {room.roomName}
                                </li></Link>
                        ))}
                        {
                            rooms.length == 0 && <CardTitle>It looks like you are not part of any room, try creating one</CardTitle>
                        }
                    </ul>
                </CardContent>
            </Card>
        </>)
}

export default Landing
