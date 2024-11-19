import { useEffect, useState } from "react";
import { Player } from "./Player";

type OtherPlayersType = {
    [key: string]: {
        username: string,
        color: string,
        x: number,
        y: number
    };
};

const GameBoard = ({ username }: { username: string }) => {
    const [otherPlayers, setOtherPlayers] = useState<OtherPlayersType>({});
    const [position, setPosition] = useState({ x: 0, y: 0 });
    const boardWidth = 1500;
    const boardHeight = 900;
    const characterSize = 40;

    const updatePositionOfOtherPlayer = (x: string, y: string, username: string, color: string) => {
        const numberX = parseInt(x) || 0
        const numberY = parseInt(y) || 0
        setOtherPlayers(prev => ({
            ...prev,
            [username]: {
                username,
                color,
                x: numberX,
                y: numberY
            }
        }))
    }

    const removeOtherPlayer = (username: string) => {
        setOtherPlayers(prev => {
            const updatedHashmap = { ...prev }
            delete updatedHashmap[username]
            return updatedHashmap
        })
    }

    useEffect(() => {
        const handleKeyDown = (e: KeyboardEvent) => {
            setPosition(prev => {
                let newX = prev.x;
                let newY = prev.y;

                switch (e.key) {
                    case 'ArrowUp':
                        newY = Math.max(0, prev.y - 10);
                        break;
                    case 'ArrowDown':
                        newY = Math.min(boardHeight - characterSize, prev.y + 10);
                        break;
                    case 'ArrowLeft':
                        newX = Math.max(0, prev.x - 10);
                        break;
                    case 'ArrowRight':
                        newX = Math.min(boardWidth - characterSize, prev.x + 10);
                        break;
                }

                return { x: newX, y: newY };
            });
        };

        window.addEventListener('keydown', handleKeyDown);
        return () => window.removeEventListener('keydown', handleKeyDown);
    }, []);

    return (
        <div
            className="relative bg-gray-200"
            style={{
                width: `${boardWidth}px`,
                height: `${boardHeight}px`
            }}
        >
            <Player x={position.x} y={position.y} color="bg-red-600" username={username} />
            {
                Object.entries(otherPlayers).map(([id, player]) => (
                    <Player key={id} y={player.y} x={player.x} color={player.color} username={player.username} />
                ))
            }
        </div>
    );
};

export default GameBoard;
