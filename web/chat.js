let ws;

function connectWebSocket() {
    ws = new WebSocket('ws://' + window.location.host + '/ws');

    ws.onopen = () => {
        console.log('Connected to WebSocket server');
        appendMessage('System', 'Connected to chat server');
    };

    ws.onmessage = (event) => {
        const response = JSON.parse(event.data);
        appendMessage('Server', response.message);
    };

    ws.onclose = () => {
        console.log('Disconnected from WebSocket server');
        appendMessage('System', 'Disconnected from chat server');
        // Try to reconnect after 5 seconds
        setTimeout(connectWebSocket, 5000);
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
        appendMessage('System', 'Error connecting to chat server');
    };
}

function appendMessage(sender, message) {
    const messagesDiv = document.getElementById('messages');
    const messageElement = document.createElement('div');
    messageElement.className = 'message';
    messageElement.textContent = `${sender}: ${message}`;
    messagesDiv.appendChild(messageElement);
    messagesDiv.scrollTop = messagesDiv.scrollHeight;
}

// Connect when the page loads
document.addEventListener('DOMContentLoaded', () => {
    connectWebSocket();
});

function sendMessage() {
    const input = document.getElementById('message-input');
    const message = input.value.trim();
    
    if (message && ws && ws.readyState === WebSocket.OPEN) {
        const data = { name: message };
        ws.send(JSON.stringify(data));
        appendMessage('You', message);
        input.value = '';
    }
}

// Handle Enter key press
document.getElementById('message-input').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') {
        sendMessage();
    }
});

// Connect when the page loads
connectWebSocket();