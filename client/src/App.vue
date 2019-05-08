<template>
  <div id="app">
    <div class="container-top">
      <div class="container-left">
        <div class="container-minimap">
          <Minimap />
        </div>
        <div class="container-targets">
          <RoomTargets />
        </div>
      </div>
      <div class="container-center">
        <div class="container-maintext">
          <MainText :windowHeight="windowHeight" />
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
import { mapState } from 'vuex';
import MainText from '@/components/MainText';
import InputBox from '@/components/InputBox';
import Minimap from '@/components/Minimap';
import RoomTargets from '@/components/RoomTargets';

export default {
  name: 'App',
  components: {
    InputBox,
    MainText,
    Minimap,
    RoomTargets
  },
  data: () => {
    return {
      windowHeight: 0,
    }
  },
  computed: mapState(['allowGlobalHotkeys']),
  methods: {
    onWindowResize() {
      this.windowHeight = window.innerHeight;
    },
    onKeyUp(event) {
      if (!this.allowGlobalHotkeys) {
        return;
      }

      let sendCommand = '';

      switch(event.key) {
        case 'w':
          sendCommand = "/move north";
          break;
        case 'a':
          sendCommand = "/move west";
          break;
        case 's':
          sendCommand = "/move south";
          break;
        case 'd':
          sendCommand = "/move east";
          break;
        case 'q':
          sendCommand = "/move down";
          break;
        case 'e':
          sendCommand = "/move up";
          break;
      }

      if (sendCommand.length > 0) {
        this.$store.dispatch('sendSlashCommand', {
          command: sendCommand
        });
      }
    }
  },
  mounted() {
    this.onWindowResize()

    window.addEventListener(
      'resize',
      this.onWindowResize
    );

    window.addEventListener(
      'keyup',
      this.onKeyUp
    );
  },
  destroyed() {
    window.removeEventListener(
      'resize',
      this.onWindowResize
    );
  }
}
</script>

<style lang="scss">
@import url('https://fonts.googleapis.com/css?family=Montserrat:100,100i,200,200i,300,300i,400,400i,500,500i,600,600i,700,700i,800,800i,900,900i');
$backgroundNormal: #111;
$backgroundLight: #1b1b1b;
$sidebarWidth: 250px;

html, body {
  padding: 0;
  margin: 0;
  height: 100%;
  background-color: $backgroundNormal;
  user-select: none;
}

::-webkit-scrollbar { width: 8px; height: 3px; }
::-webkit-scrollbar-button {  background-color: #666; }
::-webkit-scrollbar-track {  background-color: #646464; }
::-webkit-scrollbar-track-piece { background-color: #000; }
::-webkit-scrollbar-thumb { height: 50px; background-color: #666; border-radius: 3px; }
::-webkit-scrollbar-corner { background-color: #646464; }
::-webkit-resizer { background-color: #666; }

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

      .container-targets {
        flex-grow: 1;
      }
    }

    .container-center {
      flex-grow: 1;
      display: flex;
      flex-direction: column;
      box-shadow: 0px 0px 12px 1px #080808;
      position: relative;

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
