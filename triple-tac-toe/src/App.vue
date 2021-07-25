<template>
  <v-app>

    <v-main>
      <router-view/>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: 'App',

  data: () => ({
    //
  })
}
</script>

<style lang="scss">
@font-face {
  font-family: "Pixelmania";
  src: local("Pixelmania"), local("Pixelmania"),
    /*url($font_path+"Roboto-black.woff2") format("woff2"),
    url($font_path+"Roboto-black.woff") format("woff"),*/
    url("assets/fonts/Pixelmania.ttf") format("truetype");
  font-style: normal;
  font-weight: 900;
}
</style>
