import { createBrowserRouter, RouterProvider } from "react-router-dom";
import { SignUp } from "./components/Signup";
import { SignIn } from "./components/Signin";
import Landing from "./components/Landing";
import { Createroom } from "./components/CreateRoom";
import UserProvider from "./context/UserContext";

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
            <UserProvider>
                <RouterProvider router={router} />
            </UserProvider>
        </>
    );
}

export default App;
