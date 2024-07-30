const ws = new WebSocket('ws://localhost:12312/ws');
const chatbox = document.getElementById('chatbox');
const messageInput = document.getElementById('message');
const sendButton = document.getElementById('send');

ws.onmessage = function(event) {
    try {
        const message = JSON.parse(event.data);
        if (message && message.sender && message.content) {
            const formattedMessage = `${message.sender}: ${message.content}`;
            const messageElement = document.createElement('div');
            messageElement.textContent = formattedMessage;
            chatbox.appendChild(messageElement);
        } else {
            console.error("Received invalid message format:", event.data);
        }
    } catch (e) {
        console.error("Error parsing message:", e, event.data);
    }
};

sendButton.addEventListener('click', () => {
    const message = messageInput.value;
    if (message) {
        ws.send(message);
        messageInput.value = '';
    }
});

messageInput.addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        sendButton.click();
    }
});

ws.onclose = function() {
    console.log("Connection closed");
};

ws.onerror = function(error) {
    console.error("WebSocket error:", error);
};
