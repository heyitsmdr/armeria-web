<template>
  <div id="app">
    <div class="container-top">
      <div class="container-left">
        <div class="container-minimap" @click="clickMinimap">Minimap</div>
        <div class="container-targets">Room Targets</div>
      </div>
      <div class="container-center">
        <div class="container-maintext">
          <MainText />
        </div>
        <div class="container-input">
          <InputBox />
        </div>
      </div>
      <div class="container-right">Quests</div>
    </div>
    <div class="container-bottom">Bottom</div>
  </div>
</template>

<script>
import MainText from '@/components/MainText';
import InputBox from '@/components/InputBox';

export default {
  name: 'App',
  components: {
    InputBox,
    MainText
  },
  methods: {
    clickMinimap() {
      this.$store.dispatch('sendSlashCommand', {
        command: "/say hello"
      });
    }
  }
}
</script>

<style lang="scss">
@import url('https://fonts.googleapis.com/css?family=Montserrat:100,100i,200,200i,300,300i,400,400i,500,500i,600,600i,700,700i,800,800i,900,900i');
$backgroundNormal: #111;
$backgroundLight: #1b1b1b;
$sidebarWidth: 300px;

html, body {
  padding: 0;
  margin: 0;
  height: 100%;
  background-color: $backgroundNormal;
}

#app {
  font-family: 'Montserrat', sans-serif;
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
      flex-basis: $sidebarWidth;
      background-color: $backgroundLight;
      display: flex;
      flex-direction: column;

      .container-minimap {
        flex-basis: 250px;
      }
    }

    .container-center {
      flex-grow: 1;
      display: flex;
      flex-direction: column;

      .container-maintext {
        flex-grow: 1;
      }

      .container-input {
        flex-shrink: 1;
        background-color: #000;
      }
    }

    .container-right {
      flex-basis: $sidebarWidth;
      background-color: $backgroundLight;
    }
  }

  .container-bottom {
    flex-basis: 300px;
    background-color: $backgroundLight;
  }
}
</style>
