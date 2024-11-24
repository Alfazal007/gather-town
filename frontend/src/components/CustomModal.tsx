import { Button } from "@/components/ui/button"
import {
    Dialog,
    DialogContent,
    DialogHeader,
    DialogTitle,
} from "@/components/ui/dialog"

interface CustomModalProps {
    isOpen: boolean
    onClose: () => void
    onAccept: () => void
    caller: string
}

export function CustomModal({ isOpen, onClose, onAccept, caller }: CustomModalProps) {
    return (
        <Dialog open={isOpen} onOpenChange={onClose}>
            <DialogContent className="sm:max-w-[425px]">
                <DialogHeader>
                    <DialogTitle>Incoming call from {caller}. Will you leave the room for now to attend the call?</DialogTitle>
                </DialogHeader>
                <div className="grid gap-4 py-4">
                    <div className="flex justify-end space-x-2">
                        <Button variant="secondary" onClick={onClose}>
                            Reject
                        </Button>
                        <Button onClick={onAccept}>Accept</Button>
                    </div>
                </div>
            </DialogContent>
        </Dialog>
    )
}

