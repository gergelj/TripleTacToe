var socket = null
const connect = (username, vm) => {
  socket = new WebSocket('ws://' + process.env.VUE_APP_API_BASE_URL + 'ws?username=' + username)
  console.log('Attempting Connection...')

  socket.onopen = () => {
    console.log('Successfully Connected')
  }

  socket.onmessage = msg => {
    console.log(msg.data)
    const message = JSON.parse(msg.data)
    if (message.body !== '') {
      message.body = JSON.parse(message.body)
    } else {
      message.body = {}
    }
    // console.log(message)
    vm.onMessage(message)
  }

  socket.onclose = event => {
    console.log('Socket Closed Connection: ', event)
  }

  socket.onerror = error => {
    console.log('Socket Error: ', error)
    vm.onError()
  }
}

const sendMsg = msg => {
  console.log('sending msg: ', msg)
  socket.send(msg)
}

const closeConnection = () => {
  console.log('Closing connection...')
  socket.close()
}

export { connect, sendMsg, closeConnection }
