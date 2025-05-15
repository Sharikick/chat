const chat = document.getElementById("chat")
const usernameInput = document.getElementById("username")
const messageInput = document.getElementById("message")
const send = document.getElementById("btn")

const socket = new WebSocket(`ws://localhost/ws`)

socket.onmessage = function(event) {
  const msg = JSON.parse(event.data)
  const div = document.createElement('div')
  div.innerHTML = `<strong>${msg.username}:</strong> ${msg.text}`
  chat.appendChild(div)
}

send.onclick = function() {
  const username = usernameInput.value
  const message = messageInput.value

  if (username && message) {
    socket.send(JSON.stringify({ username, text: message }));
    messageInput.value = '';
  } else {
    alert("Введите имя и сообщение")
  }
}
