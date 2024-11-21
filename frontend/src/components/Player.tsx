export const Player = ({ x, y, color, username }: { x: number, y: number, color: string, username: string }) => (
    <div
        className={`absolute ${color}`}
        style={{
            left: `${x}px`,
            top: `${y}px`,
            width: "100px",
            height: "40px",
            borderRadius: "5px",
            paddingLeft: "4px",
            paddingTop: "3px"
        }}>
        {username?.substring(0, Math.min(username.length, 10))}
    </div>
);

