<template>
  <div>
    <v-btn v-bind="size" @click="buttonClick" btn fab tile elevation="0" :outlined="disableClick" :dark="!disabled" :color="getColor">
      <v-icon>{{getIcon}}</v-icon>
    </v-btn>
  </div>
</template>
<script>
export default {
  name: 'Badge',
  props: {
    number: {
      required: true,
      type: Number
    },
    disabled: {
      type: Boolean,
      default: false
    },
    position: {
      type: Array
    },
    player: {
      type: Number
    },
    currentPlayer: {
      type: Number
    }
  },
  methods: {
    buttonClick () {
      if (!this.disableClick && this.number === -1) {
        this.$emit('click', this.position)
      }
    }
  },
  computed: {
    getColor () {
      switch (this.number) {
        case 0: return 'error' // O
        case 1: return 'primary' // X
        case -1: return 'success' // EMPTY
        case 2: return 'deep-purple' // DRAW
      }
      return ''
    },
    getIcon () {
      switch (this.number) {
        case 0: return 'mdi-circle-outline'
        case 1: return 'mdi-close'
        case -1: return ''
        case 2: return 'mdi-help'
      }
      return ''
    },
    myTurn () {
      return this.player === this.currentPlayer
    },
    disableClick () {
      return this.disabled || !this.myTurn
    },
    size () {
      const breakpoint = this.$vuetify.breakpoint.name
      const order = { xs: 'x-small', sm: 'x-small', md: 'small', large: 'large', xl: 'large' }
      // const size = ['xs': 'x-small', 'sm': 'small', 'lg': 'large', 'xl': 'x-large']
      return order[breakpoint] ? { [order[breakpoint]]: true } : {}
    }
  }
}
</script>
