const socket = new WebSocket('ws://localhost/ws/');

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

  switch (messageType) {
    case 0:
      const username = data.slice(4, data.length-1)
      element = createUser(view.getUint8(data.length-1), decoder.decode(username))
      updatePosition(element, view.getUint16(0, false), view.getUint16(2, false))
      break
    case 1:
      element = document.getElementById("user-" + view.getUint8(4))
      updatePosition(element, view.getUint16(0, false), view.getUint16(2, false))
      break
    case 2:
      deleteUser(view.getUint8(0))
      break
    case 3:
      updatePosition(null, view.getUint16(0, false), view.getUint16(2, false))
      break
    case 4:
      addRequest(view.getUint8(1, false), view.getUint8(0, false), "received")
      break
    case 5:
      if (view.getUint8(0, false) == 1) {
        addRequest(view.getUint8(1, false), null, "sent")
      }
      if (view.getUint8(0, false) == 2) {
        const gameID = view.getUint8(2, false)

        const buffer = new ArrayBuffer(1);
        new DataView(buffer).setUint8(0, 6)
        sendMessage(buffer)

        window.location.replace(window.location.origin + "/game/" + gameID);
      }
      break
  }
};

socket.onclose = function(event) {

};

function sendCoords(x, y) {
  const buffer = new ArrayBuffer(5);
  const view = new DataView(buffer);

  view.setUint8(0, 1);
  view.setUint16(1, x, false);
  view.setUint16(3, y, false);

  socket.send(buffer);
}
function sendMessage(message) {
  socket.send(message);
}