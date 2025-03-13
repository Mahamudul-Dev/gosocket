let socket = new WebSocket('ws://localhost:8000/ws');
socket.onopen = function(ws, event) {
    console.log("Connection established!");

    socket.send("Hello from client");

    
socket.onmessage = (event) => {
  console.log("received from server", event.data);
};

socket.onclose = function (event) {
  if (event.wasClean) {
    console.log(
      "Connection closed clean, code=" + event.code + " reason=" + event.reason
    );
  } else {
    // e.g. server process killed or network down
    // event.code is usually 1006 in this case
    console.log("Connection died");
  }
};
};
