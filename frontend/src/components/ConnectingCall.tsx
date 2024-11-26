import { Card } from "@/components/ui/card"
import { Phone } from 'lucide-react'

export default function ConnectingToCall() {
    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100">
            <Card className="w-full max-w-md p-6 space-y-6">
                <div className="flex flex-col items-center space-y-4">
                    <div className="relative">
                        <div className="w-24 h-24 border-t-4 border-blue-500 border-solid rounded-full animate-spin"></div>
                        <Phone className="w-12 h-12 text-blue-500 absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2" />
                    </div>
                    <h1 className="text-2xl font-bold text-center">Connecting to Call</h1>
                    <p className="text-gray-600 text-center">Please wait while we connect you to the call...</p>
                </div>
            </Card>
        </div>
    )
}

