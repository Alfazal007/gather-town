export const Player = ({ x, y, color, username }: { x: number, y: number, color: string, username: string }) => (
    <div
        className={`absolute ${color}`}
        style={{
            left: `${x}px`,
            top: `${y}px`,
            width: "40px",
            height: "40px"
        }}>
        {username}
    </div>
);

