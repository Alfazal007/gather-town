import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";

const WS_URL = "ws://localhost:8001/ws";

export const useSocket = () => {
    const [socket, setSocket] = useState<WebSocket | null>(null);
    const navigate = useNavigate()
    useEffect(() => {
        const ws = new WebSocket(WS_URL);
        ws.onopen = () => {
            setSocket(ws);
        };
        ws.onclose = () => {
            setSocket(null);
        };
        return () => {
            //            ws.close();
            navigate("/")
        };
    }, []);
    return socket;
};

