<template>
  <h1 id="tripletactoe" class="text-center">TRIPLE=TAC=TOE</h1>
</template>

<script>
export default {
  name: 'Logo'
}
</script>

<style>
  #tripletactoe {
    font-family: "Pixelmania", Helvetica, Arial;
    background-image: linear-gradient(90deg, #ff5252, #1976d2);
    display: inline-block;
    background-clip: text;
    -webkit-background-clip: text;
    -moz-background-clip: text;
    -webkit-text-fill-color: transparent;
    -moz-text-fill-color: transparent;
  }
  @media (max-width: 550px) {
    #tripletactoe { font-size: 1em; }
  }

  @media (min-width: 551px) {
    #tripletactoe { font-size: 1.5rem; }
  }
</style>
