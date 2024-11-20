import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { SignUp } from "./components/Signup";
import { SignIn } from "./components/Signin";
import Landing from "./components/Landing";
import { Createroom } from "./components/CreateRoom";
import UserProvider from "./context/UserContext";
import { ThemeProvider } from "./components/theme-provider";
import RoomOuter from "./components/RoomOuter";

export interface User {
    accessToken: string;
    refreshToken: string;
    username: string;
    id: string;
}

function App() {

    const router = createBrowserRouter([
        {
            path: "/signup",
            element: <SignUp />,
        },
        {
            path: "/signin",
            element: <SignIn />,
        },
        {
            path: "/",
            element: <Landing />
        },
        {
            path: "/create-room",
            element: <Createroom />
        },
        {
            path: "/room/:roomId",
            element: <RoomOuter />,
        }

    ]);

    return (
        <>
            <ThemeProvider defaultTheme="dark" storageKey="vite-ui-theme">
                <UserProvider>
                    <RouterProvider router={router} />
                </UserProvider>
            </ThemeProvider>
        </>
    );
}

export default App;
