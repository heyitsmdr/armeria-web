<template>
  <div id="app">
    <div class="container-top">
      <div class="container-left">Top Left</div>
      <div class="container-center">Top Center</div>
      <div class="container-right">Top Right</div>
    </div>
    <div class="container-bottom">Bottom</div>
  </div>
</template>

<script>

export default {
  name: 'app',
  components: {
    
  }
}
</script>

<style lang="scss">
html, body {
  padding: 0;
  margin: 0;
  height: 100%;
  background-color: #111;
}

#app {
  font-family: 'Comfortaa', cursive;
  /* -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale; */
  font-size: 14px;
  margin: 0;
  padding: 0;
  color: #dbe3e6;
  display: flex;
  flex-direction: column;
  height: 100%;

  .container-top {
    flex-grow: 1;
    display: flex;

    .container-left {
      flex-basis: 250px;
    }

    .container-center {
      flex-grow: 1;
    }

    .container-right {
      flex-basis: 250px;
    }
  }

  .container-bottom {
    flex-basis: 300px;
  }
}
</style>
