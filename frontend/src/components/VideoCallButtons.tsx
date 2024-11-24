import { useEffect, useState } from "react"
import { Button } from "@/components/ui/button"
import { ScrollArea } from "@/components/ui/scroll-area"
import { Video } from 'lucide-react'
import { OtherPlayersType } from "./GameBoard"

export default function VideoCall({ otherPeople, makeACall, acceptor, disconnectAndMoveToCallAttendScreen }: { otherPeople: OtherPlayersType, makeACall: (receiver: string) => void, acceptor: string, disconnectAndMoveToCallAttendScreen: (username: string) => void }) {
    const [selectedUser, setSelectedUser] = useState<string | null>(null)

    function initiateCall(username: string) {
        setSelectedUser(username)
        makeACall(username)
    }

    useEffect(() => {
        if (acceptor == selectedUser) {
            console.log("Call has been accepted")
            disconnectAndMoveToCallAttendScreen(selectedUser)
        }
    }, [acceptor])

    return (
        <>
            <div className="w-1/5 h-80 fixed bottom-2 right-4 mr-2 flex flex-col bg-gray-50 dark:bg-gray-800 border-r border-t border-gray-200 dark:border-gray-700 rounded-tr-lg shadow-lg">
                <div className="bg-primary text-primary-foreground p-2 flex justify-between items-center">
                    <h2 className="text-sm font-semibold">Talk to</h2>
                </div>
                <ScrollArea className="flex-grow">
                    <div className="space-y-6 p-4">
                        {Object.entries(otherPeople).map(([id, player]) => (
                            <Button
                                key={id}
                                variant="outline"
                                size="lg"
                                onClick={() => initiateCall(player.username)}
                                className="w-full py-6 text-lg"
                            >
                                {player.username}
                                <Video className="ml-2 h-6 w-6" />
                            </Button>
                        ))}
                    </div>
                </ScrollArea>
            </div>
        </>
    )
}

