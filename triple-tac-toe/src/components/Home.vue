<template>
  <div class="ml-5 mr-5 mt-5">
    <v-dialog v-model="showDialog" width="300" persistent>
      <v-card>
        <v-card-title v-if="status === 'DRAW'">It's a draw. Play again?</v-card-title>
        <v-card-title v-else-if="status === 'YOU_WON'">You won! Play again?</v-card-title>
        <v-card-title v-else-if="status === 'OPPONENT_WON'">Opponent won. Play again?</v-card-title>
        <v-card-title v-else-if="status === 'USERNAME_TAKEN'">Name is taken <br> Choose another one</v-card-title>
        <v-card-title v-else-if="status === 'LEFT'">Opponent has left. <br> Play again?</v-card-title>
        <v-card-title v-else>Enter you name to play</v-card-title>
        <v-card-text>
          <v-text-field label="Name" v-model="name"></v-text-field>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn @click="connect" color="success" :disabled="name === ''">Play</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>

    <v-card flat max-width="850" class="mx-auto">
      <v-container fluid>
        <v-row class="justify-center">
          <logo></logo>
        </v-row>
      </v-container>
    </v-card>

    <v-container>
      <v-card flat max-width="850" class="mx-auto">
        <v-card-title class="d-flex justify-center">
          <div class="monotext" v-if="status === 'WAIT'">Waiting for opponent</div>
          <div class="monotext" v-if="status === 'PLAY'">{{player === currentPlayer ? 'My turn' : 'Opponent\'s turn'}}</div>
        </v-card-title>
      </v-card>
      <v-card flat max-width="850" class="mx-auto pa-0">
        <toolbar v-if="status === 'PLAY'" @leave="leave" :name="name" :player="player" :currentPlayer="currentPlayer" :opponentName="opponentName"></toolbar>
        <v-card-text v-if="showGameBoard" class="pa-0 ma-0">
          <v-container class="ma-0">
            <v-row class="justify-center flex-nowrap" v-for="i in [0,1,2]" :key="i">
              <!--<v-col class="pa-1" v-for="j in [0,1,2]" :key="j" cols="4">-->
                <v-card class="ma-1" v-for="j in [0,1,2]" :key="j" :style="(bigGrid[i][j].allowed ? 'background-color: rgb(234, 255, 234) !important;' : '')" outlined>
                  <v-container fluid>
                    <v-row class="justify-space-around flex-nowrap" v-for="k in [0,1,2]" :key="k">
                      <!--<v-col class="pa-1" v-for="l in [0,1,2]" :key="l">-->
                        <badge class="ma-1" v-for="l in [0,1,2]" :key="l" @click="clicked" :number="grid[i][j][k][l]" :disabled="!bigGrid[i][j].allowed" :player="player" :currentPlayer="currentPlayer" :position="[i, j, k, l]"></badge>
                      <!--</v-col>-->
                    </v-row>
                  </v-container>
                </v-card>
              <!--</v-col>-->
            </v-row>
          </v-container>
        </v-card-text>
      </v-card>
    </v-container>
  </div>
</template>

<script>
import { connect, sendMsg, closeConnection } from '@/websocket/websocket.js'
import Badge from './Badge.vue'
import Toolbar from './Toolbar.vue'
import Logo from './Logo.vue'
export default {
  components: { Badge, Toolbar, Logo },
  name: 'Home',
  data () {
    return {
      showDialog: true, // change to true!!!
      name: '',
      opponentName: '',
      status: 'CREATED',
      player: -1, // change to -1
      currentPlayer: -1, // change to -1
      grid: {},
      bigGrid: {},
      keepAlive: 20 // seconds
    }
  },
  created () {
    this.initializeGame()
  },
  mounted () {
    this.name = localStorage.getItem('name')
    setInterval(() => {
      if (this.status === 'WAIT' || this.status === 'PLAY') {
        sendMsg(JSON.stringify({ type: 'KEEP_ALIVE' }))
      }
    }, this.keepAlive * 1000)
  },
  methods: {
    connect () {
      localStorage.setItem('name', this.name)
      connect(this.name, this)
    },
    close () {
      closeConnection()
    },
    leave () {
      this.status = 'CREATED'
      this.resetGame()
    },
    resetGame () {
      this.close()
      this.showDialog = true
      this.opponentName = ''
      this.player = -1
      this.currentPlayer = -1
      this.grid = {}
      this.bigGrid = {}
      this.initializeGame()
    },
    onMessage (message) {
      switch (message.type) {
        case 'WAIT': this.onWaitMessage(); break
        case 'START_GAME': this.onStartGameMessage(message); break
        case 'TURN': this.onTurnMessage(message); break
        case 'USERNAME_TAKEN': this.onUsernameTakenMessage(); break
        case 'LEFT': this.onLeftMessage(); break
      }
      console.log(message)
    },
    onError () {
      this.resetGame()
    },
    onUsernameTakenMessage () {
      this.status = 'USERNAME_TAKEN'
      // this.close() // maybe server closes connection
    },
    onLeftMessage () {
      this.status = 'LEFT'
      this.resetGame()
    },
    onWaitMessage () {
      this.showDialog = false
      this.status = 'WAIT'
    },
    onStartGameMessage (message) {
      this.showDialog = false
      this.player = message.body.number
      this.currentPlayer = message.body.startNumber
      this.opponentName = message.body.opponent
      this.status = 'PLAY'
    },
    onTurnMessage (message) {
      this.currentPlayerStep(message.body)
    },
    initializeGame () {
      for (let i = 0; i < 3; i++) {
        this.grid[i] = {}
        this.bigGrid[i] = {}
        for (let j = 0; j < 3; j++) {
          this.grid[i][j] = {}
          this.bigGrid[i][j] = { value: -1, allowed: true }
          for (let k = 0; k < 3; k++) {
            this.grid[i][j][k] = {}
            for (let l = 0; l < 3; l++) {
              this.grid[i][j][k][l] = -1
            }
          }
        }
      }
    },
    clicked (position) {
      const message = { type: 'TURN', body: position }
      sendMsg(JSON.stringify(message))
      this.currentPlayerStep(position)
    },
    currentPlayerStep (position) {
      let i, j, k, l
      // eslint-disable-next-line prefer-const
      [i, j, k, l] = position
      this.grid[i][j][k][l] = this.currentPlayer
      this.checkGrid(i, j, k, l)
      this.changePlayer()
      this.setAllowedTiles(k, l)
    },
    checkGrid (i, j, k, l) {
      if (this.checkSmallGrid(i, j)) {
        this.playerWonSmallGrid(i, j)
        if (this.checkBigGrid()) {
          this.gameOver(this.currentPlayer)
        } else if (this.checkBigTie()) {
          this.gameOver(2)
        }
      } else if (this.checkSmallTie(i, j)) {
        for (let x = 0; x < 3; x++) {
          for (let y = 0; y < 3; y++) {
            this.grid[i][j][x][y] = 2
          }
        }
        this.bigGrid[i][j].value = 2
        if (this.checkBigTie()) {
          this.gameOver(2)
        }
      }
    },
    checkSmallGrid (x, y) {
      return (this.grid[x][y][0][0] === this.currentPlayer && this.grid[x][y][0][1] === this.currentPlayer && this.grid[x][y][0][2] === this.currentPlayer) ||
                (this.grid[x][y][1][0] === this.currentPlayer && this.grid[x][y][1][1] === this.currentPlayer && this.grid[x][y][1][2] === this.currentPlayer) ||
                (this.grid[x][y][2][0] === this.currentPlayer && this.grid[x][y][2][1] === this.currentPlayer && this.grid[x][y][2][2] === this.currentPlayer) ||
                (this.grid[x][y][0][0] === this.currentPlayer && this.grid[x][y][1][0] === this.currentPlayer && this.grid[x][y][2][0] === this.currentPlayer) ||
                (this.grid[x][y][0][1] === this.currentPlayer && this.grid[x][y][1][1] === this.currentPlayer && this.grid[x][y][2][1] === this.currentPlayer) ||
                (this.grid[x][y][0][2] === this.currentPlayer && this.grid[x][y][1][2] === this.currentPlayer && this.grid[x][y][2][2] === this.currentPlayer) ||
                (this.grid[x][y][0][0] === this.currentPlayer && this.grid[x][y][1][1] === this.currentPlayer && this.grid[x][y][2][2] === this.currentPlayer) ||
                (this.grid[x][y][2][0] === this.currentPlayer && this.grid[x][y][1][1] === this.currentPlayer && this.grid[x][y][0][2] === this.currentPlayer)
    },
    checkSmallTie (x, y) {
      for (let i = 0; i < 3; i++) {
        for (let j = 0; j < 3; j++) {
          if (this.grid[x][y][i][j] === -1) {
            return false
          }
        }
      }
      return true
    },
    checkBigGrid () {
      return (this.bigGrid[0][0].value === this.currentPlayer && this.bigGrid[0][1].value === this.currentPlayer && this.bigGrid[0][2].value === this.currentPlayer) ||
                (this.bigGrid[1][0].value === this.currentPlayer && this.bigGrid[1][1].value === this.currentPlayer && this.bigGrid[1][2].value === this.currentPlayer) ||
                (this.bigGrid[2][0].value === this.currentPlayer && this.bigGrid[2][1].value === this.currentPlayer && this.bigGrid[2][2].value === this.currentPlayer) ||
                (this.bigGrid[0][0].value === this.currentPlayer && this.bigGrid[1][0].value === this.currentPlayer && this.bigGrid[2][0].value === this.currentPlayer) ||
                (this.bigGrid[0][1].value === this.currentPlayer && this.bigGrid[1][1].value === this.currentPlayer && this.bigGrid[2][1].value === this.currentPlayer) ||
                (this.bigGrid[0][2].value === this.currentPlayer && this.bigGrid[1][2].value === this.currentPlayer && this.bigGrid[2][2].value === this.currentPlayer) ||
                (this.bigGrid[0][0].value === this.currentPlayer && this.bigGrid[1][1].value === this.currentPlayer && this.bigGrid[2][2].value === this.currentPlayer) ||
                (this.bigGrid[2][0].value === this.currentPlayer && this.bigGrid[1][1].value === this.currentPlayer && this.bigGrid[0][2].value === this.currentPlayer)
    },
    checkBigTie () {
      for (let i = 0; i < 3; i++) {
        for (let j = 0; j < 3; j++) {
          if (this.bigGrid[i][j].value === -1) {
            return false
          }
        }
      }
      return true
    },
    playerWonSmallGrid (x, y) {
      this.bigGrid[x][y].value = this.currentPlayer
      for (let i = 0; i < 3; i++) {
        for (let j = 0; j < 3; j++) {
          this.grid[x][y][i][j] = this.currentPlayer
        }
      }
    },
    gameOver (code) {
      if (code === 2) {
        this.status = 'DRAW'
      } else {
        if (code === this.player) {
          this.status = 'YOU_WON'
        } else {
          this.status = 'OPPONENT_WON'
        }
      }
    },
    setAllowedTiles (k, l) {
      if (this.bigGrid[k][l].value !== -1) { // big grid tile is not empty then allow all empty tiles
        for (let i = 0; i < 3; i++) {
          for (let j = 0; j < 3; j++) {
            if (this.bigGrid[i][j].value === -1) {
              this.bigGrid[i][j].allowed = true
            } else {
              this.bigGrid[i][j].allowed = false
            }
          }
        }
      } else {
        for (let i = 0; i < 3; i++) {
          for (let j = 0; j < 3; j++) {
            if (i === k && j === l) {
              this.bigGrid[i][j].allowed = true
            } else {
              this.bigGrid[i][j].allowed = false
            }
          }
        }
      }
    },
    changePlayer () {
      this.currentPlayer = this.currentPlayer === 0 ? 1 : 0
    },
    toIndex (i, j, k, l) {
      return 27 * i + 9 * j + 3 * k + l
    },
    toIndex2 (i, j) {
      return 3 * i + j
    }
  },
  computed: {
    opponent () {
      if (this.player === 0) {
        return 1
      } else if (this.player === 1) {
        return 0
      } else {
        return -1
      }
    },
    showGameBoard () {
      return this.status === 'PLAY'
    }
  },
  watch: {
    status (newStatus) {
      switch (newStatus) {
        case 'DRAW':
        case 'YOU_WON':
        case 'OPPONENT_WON': {
          this.resetGame()
        }
      }
    }
  }
}
</script>

<style>
  .monotext {
    font-family: 'Courier New', monospace;
    font-weight: 600;
  }
</style>
