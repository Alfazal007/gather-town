import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { SignUp } from "./components/Signup";
import { SignIn } from "./components/Signin";
import Landing from "./components/Landing";
import { Createroom } from "./components/CreateRoom";
import UserProvider from "./context/UserContext";
import { ThemeProvider } from "./components/theme-provider";
import Navbar from "./components/Navbar";

export interface User {
    accessToken: string;
    refreshToken: string;
    username: string;
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
