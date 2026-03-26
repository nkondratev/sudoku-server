let socket = new WebSocket("ws://localhost:8080/ws")

  
socket.onmessage = (e)=> {
  console.log(e.data)
}
