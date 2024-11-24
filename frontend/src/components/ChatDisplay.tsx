import { ScrollArea } from "@/components/ui/scroll-area"
import { Button } from "./ui/button"
import { Send, X } from "lucide-react"
import { useState, useRef, useEffect } from "react"
import { Input } from "./ui/input"
import VideoCall from "./VideoCallButtons"
import { nanoid } from "nanoid"
import { OtherPlayersType } from "./GameBoard"

interface Message {
    sender: string
    message: string
}

interface ChatDisplayProps {
    messages: Message[],
    onLeaveRoom: () => void,
    onSendMessage: (message: string) => void,
    otherPeople: OtherPlayersType
}

export default function ChatDisplay({ messages, onLeaveRoom, onSendMessage, otherPeople }: ChatDisplayProps) {
    const [newMessage, setNewMessage] = useState("")
    const messagesEndRef = useRef<HTMLDivElement>(null)

    const scrollToBottom = () => {
        messagesEndRef.current?.scrollIntoView({ behavior: "smooth" })
    }

    const handleSendMessage = (e: React.FormEvent) => {
        e.preventDefault()
        if (newMessage.trim()) {
            onSendMessage(newMessage)
            setNewMessage('')
        }
    }

    useEffect(() => {
        scrollToBottom()
    }, [messages])

    return (
        <>
            <div className="w-1/5 h-1/2 fixed top-2 right-4 mr-2 top-1/4 flex flex-col bg-gray-50 dark:bg-gray-800 border-l border-t border-gray-200 dark:border-gray-700 rounded-tl-lg shadow-lg">
                <div className="bg-primary text-primary-foreground p-2 flex justify-between items-center rounded-tl-lg">
                    <h2 className="text-sm font-semibold">Chat Messages</h2>
                    <Button
                        variant="ghost"
                        size="icon"
                        onClick={onLeaveRoom}
                        aria-label="Leave Room"
                        className="h-6 w-6"
                    >
                        <X className="h-4 w-4" />
                    </Button>
                </div>
                <ScrollArea className="flex-grow overflow-y-auto">
                    <div className="p-2 space-y-2">
                        {messages.map((message) => (
                            <div key={nanoid()} className="text-sm">
                                <div className="font-semibold text-xs text-muted-foreground">
                                    {message.sender}
                                </div>
                                <div className="bg-white dark:bg-gray-700 p-2 rounded-md mt-1 shadow-sm">
                                    {message.message}
                                </div>
                            </div>
                        ))}
                        <div ref={messagesEndRef} />
                    </div>
                </ScrollArea>
                <form onSubmit={handleSendMessage} className="p-2 border-t border-gray-200 dark:border-gray-700">
                    <div className="flex items-center space-x-2">
                        <Input
                            type="text"
                            placeholder="Type a message..."
                            value={newMessage}
                            onChange={(e) => setNewMessage(e.target.value)}
                            className="flex-grow"
                        />
                        <Button type="submit" size="icon" aria-label="Send message">
                            <Send className="h-4 w-4" />
                        </Button>
                    </div>
                </form>
            </div>
            <VideoCall otherPeople={otherPeople} />
        </>
    )
}

