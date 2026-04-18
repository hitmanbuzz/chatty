interface Message {
    user_id: number;
    content: string;
}

const SERVER_URL = "ws://localhost:6769/echo"

const socket = new WebSocket(SERVER_URL);

socket.addEventListener("open", () => {
    console.log("Conected to websocket server");

    const payload: Message = {
        user_id: 1,
        content: "Hello from js bun client",
    };

    socket.send(JSON.stringify(payload))
    console.log("Payload Sent:", payload)
})

socket.addEventListener("message", (event) => {
    console.log("Raw data received:", event.data);

    try {
        const incomingData = JSON.parse(event.data) as Message;
        console.log(`Message from ${incomingData.user_id}: ${incomingData.content}`);
    } catch (error) {
        console.error("Could not parse incoming message as JSON");
    }
});

socket.addEventListener("error", (event) => {
    console.error("WebSocket Error occurred:", event);
});

socket.addEventListener("close", (event) => {
    console.log(`Connection closed. Code: ${event.code}, Reason: ${event.reason || "None"}`);
});
