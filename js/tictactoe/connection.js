const socket = new WebSocket('ws://localhost/ws' + window.location.pathname);

socket.binaryType = "arraybuffer";
const decoder = new TextDecoder();

socket.onopen = function(event) {
};

socket.onmessage = function(event) {
    const buffer = event.data;
    const message = new Uint8Array(buffer);

    const messageType = message[0]
    const data = message.slice(1)
    const view = new DataView(data.buffer)

    console.log(message)

    switch (messageType) {
        case 0:
            loadOpponent(view.getUint8(0), decoder.decode(data.slice(1)))
            break
        case 7:
            opponentMove(view.getUint8(0))
            break
        case 8:
            gameUpdate(view.getUint8(0), data.slice(1))
            break
    }
};

socket.onclose = function(event) {

};

function sendMove(move) {
    const buffer = new ArrayBuffer(3);
    const view = new DataView(buffer);

    view.setUint8(0, 7)
    view.setUint8(1, move)
    view.setUint8(2, gameID)

    sendMessage(buffer);
}

function sendMessage(message) {
    socket.send(message);
}