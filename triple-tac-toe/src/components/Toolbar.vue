<template>
    <v-card-title class="d-flex flex-nowrap">
      <div class="d-flex flex-row">
        <v-icon v-if="currentPlayer === player" x-large color="green">mdi-chevron-right</v-icon>
        <div class="d-flex flex-column align-center flex-sm-column flex-md-row flex-lg-row flex-xl-row">
          <badge class="mr-5" :number="player"></badge>
          <h3 class="monotext">{{name}}</h3>
        </div>
      </div>

      <v-spacer></v-spacer>
      <v-btn @click="$emit('leave')" plain>Leave</v-btn>
      <v-spacer></v-spacer>

      <div class="d-flex flex-row">
        <div class="d-flex flex-column-reverse align-center flex-sm-column-reverse flex-md-row flex-lg-row flex-xl-row">
          <h3 class="monotext">{{opponentName}}</h3>
          <badge class="ml-5" :number="opponent"></badge>
        </div>
        <v-icon v-if="currentPlayer === opponent" x-large color="green">mdi-chevron-left</v-icon>
      </div>
    </v-card-title>
</template>

<script>
import Badge from './Badge.vue'
export default {
  components: { Badge },
  name: 'Toolbar',
  props: {
    name: {
      required: true,
      type: String
    },
    player: {
      required: true,
      type: Number
    },
    currentPlayer: {
      required: true,
      type: Number
    },
    opponentName: {
      required: true,
      type: String
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
    }
  }
}
</script>
