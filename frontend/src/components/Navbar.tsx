import { Home, PlusCircle, LogOut } from 'lucide-react'
import { Button } from "@/components/ui/button"
import { Link, useNavigate } from 'react-router-dom'
import { ModeToggle } from './mode-toggle'
import Cookies from 'js-cookie';
import { User } from '@/App'
import { useToast } from '@/hooks/use-toast';

const deleteCookie = (name: string) => {
    Cookies.remove(name);
};

export default function Navbar({ setUser }: {
    setUser: React.Dispatch<React.SetStateAction<User | null>>
}) {
    const { toast } = useToast();
    const navigate = useNavigate()
    const handleLogout = async () => {
        setUser(null)
        toast({
            title: "Logged out successfully",
        })
        deleteCookie("accessToken")
        deleteCookie("refreshToken")
        navigate('/signin')
    }

    return (
        <nav className="bg-background shadow-md">
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                <div className="flex justify-between h-16">
                    <div className="flex">
                        <div className="flex-shrink-0 flex items-center">
                            <span className="text-xl font-bold text-foreground">Gather Town</span>
                        </div>
                        <div className="hidden sm:ml-6 sm:flex sm:space-x-8">
                            <Link
                                to="/"
                                className="inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium text-muted-foreground hover:border-border hover:text-foreground"
                            >
                                <Home className="mr-1 h-5 w-5" />
                                Home
                            </Link>
                            <Link
                                to="/create-room"
                                className="inline-flex items-center px-1 pt-1 border-b-2 border-transparent text-sm font-medium text-muted-foreground hover:border-border hover:text-foreground"
                            >
                                <PlusCircle className="mr-1 h-5 w-5" />
                                Create Room
                            </Link>
                        </div>
                        <div className="flex items-center">
                            <div
                                className="inline-flex items-center px-10 py-2 border border-transparent text-sm font-medium rounded-md text-muted-foreground hover:text-foreground focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
                            >
                                <ModeToggle />
                            </div>
                        </div>
                    </div>

                    <div className="flex items-center">
                        <Button
                            variant="ghost"
                            onClick={handleLogout}
                            className="inline-flex items-center px-3 py-2 border border-transparent text-sm font-medium rounded-md text-muted-foreground hover:text-foreground focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
                        >
                            <LogOut className="mr-2 h-5 w-5" />
                            Logout
                        </Button>
                    </div>
                </div>
            </div>
        </nav>
    )
}
