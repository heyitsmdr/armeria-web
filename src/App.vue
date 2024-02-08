<template>
  <div id="app">
    <div class="container-wrapper">
      <div class="container-left" :style="{ display: leftSidebar }">
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
        <div class="container-bars">
          <Vitals />
        </div>
      </div>
      <div class="container-right" :style="{ display: rightSidebar }">
        <div class="container-skills">
          <Skills />
        </div>
        <div class="container-inventory">
          <Inventory />
        </div>
      </div>
    </div>
    <div class="status-bar-container">
      <StatusBar />
    </div>
    <ItemTooltip />
    <ContextMenu />
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from "vue";
import { useStore } from "vuex";
import InputBox from "@/components/InputBox.vue";
import MainText from "@/components/MainText.vue";
import Minimap from "@/components/Minimap.vue";
import RoomTargets from "@/components/RoomTargets.vue";
import Inventory from "@/components/Inventory.vue";
import Vitals from "@/components/Vitals.vue";
import Skills from "@/components/Skills.vue";
import StatusBar from "@/components/StatusBar.vue";
import ItemTooltip from "@/components/ItemTooltip.vue";
import ContextMenu from "@/components/ContextMenu.vue";

const windowHeight = ref(0);
const windowWidth = ref(0);
const leftSidebar = ref("flex");
const rightSidebar = ref("flex");

const store = useStore();
const allowGlobalHotkeys = computed(() => store.state.allowGlobalHotkeys);
const isConnected = computed(() => store.state.isConnected);
const playerInfo = computed(() => store.state.playerInfo);
const contextMenuVisible = computed(() => store.state.contextMenu.visible);
const hasPermission = computed(() => store.getters.hasPermission);

watch(isConnected, (connected) => {
  let token = store.state.autoLoginToken;
  if (connected) {
    store.dispatch("showText", { data: `Welcome to Armeria!\n\n` });

    if (token.length > 0) {
      const char = token.split(":")[0];
      store.dispatch("showText", {
        data: `You are automatically being logged in as '${char}'.\n`,
      });
      store.dispatch("sendSlashCommand", {
        command: `/logintoken ${token}`,
        hidden: true,
      });
    } else {
      store.dispatch("showText", {
        data: "If you have an existing character, you can <b>/login</b>. Otherwise, <b>/create</b> a new one.\n",
      });
    }

    // This "keep alive" is needed for Heroku. Otherwise, if the socket
    // connection is idle for 55 seconds, the Heroku load balancer will
    // terminate the connection and throw an H15 error.
    window.socketKeepAlive = setInterval(sendKeepAlive, 20000);
    sendKeepAlive();
  } else {
    window.document.title = "**DISCONNECTED** Armeria.io";
    clearInterval(window.socketKeepAlive);
    window.socketKeepAlive = null;
  }
});

watch(playerInfo, (info) => {
  window.document.title = `${info.name} - Armeria.io`;
});

onMounted(() => {
  onWindowResize();
  window.addEventListener("resize", onWindowResize);
  window.addEventListener("keyup", onKeyUp);
});

function sendKeepAlive() {
  store.dispatch("sendKeepAlive");
}

function onWindowResize() {
  windowHeight.value = window.innerHeight;
  windowWidth.value = window.innerWidth;
  if (windowWidth.value < 784) {
    leftSidebar.value = "none";
  } else {
    leftSidebar.value = "flex";
  }

  if (windowWidth.value < 1035) {
    // showRightSidebar.value = false;
    rightSidebar.value = "none";
  } else {
    // showRightSidebar.value = true;
    rightSidebar.value = "flex";
  }
  //document.querySelector('.container-center').style.maxWidth = `${windowWidth.value-500}px`;

  // If the context menu is open, let's hide it since the window is being resized.
  if (contextMenuVisible.value) {
    store.dispatch("contextMenu/hide");
  }
}

function onKeyUp(event) {
  if (!allowGlobalHotkeys.value) {
    return;
  }

  let moveCommand = "";

  switch (event.key.toLowerCase()) {
    case "w":
      moveCommand = "/move north";
      break;
    case "a":
      moveCommand = "/move west";
      break;
    case "s":
      moveCommand = "/move south";
      break;
    case "d":
      moveCommand = "/move east";
      break;
    case "q":
      moveCommand = "/move down";
      break;
    case "e":
      moveCommand = "/move up";
      break;
    case "escape":
      store.dispatch("setObjectTarget", "");
      break;
    case "enter":
      store.dispatch("setForceInputFocus", { forced: true });
      break;
    case "/":
      store.dispatch("setForceInputFocus", { forced: true, text: "/" });
      break;
  }

  if (moveCommand.length > 0) {
    if (hasPermission.value("CAN_BUILD")) {
      if (event.shiftKey) {
        moveCommand = moveCommand.replace("/move", "/room create");
      } else if (event.ctrlKey) {
        moveCommand = moveCommand.replace("/move", "/room destroy");
      }
    }

    store.dispatch("sendSlashCommand", {
      command: moveCommand,
      hidden: true,
    });
  }
}
</script>

<style lang="scss">
@import "@/styles/common";
$backgroundNormal: #111;
$backgroundLight: #1b1b1b;
$sidebarWidth: 250px;

html,
body {
  padding: 0;
  margin: 0;
  height: 100%;
  background-color: $bg-color;
  user-select: none;
}

::-webkit-scrollbar {
  width: 3px;
  height: 3px;
}
::-webkit-scrollbar-button {
  background-color: #666;
}
::-webkit-scrollbar-track {
  background-color: #646464;
}
::-webkit-scrollbar-track-piece {
  background-color: #111;
}
::-webkit-scrollbar-thumb {
  height: 50px;
  background-color: #333;
  border-radius: 0px;
}
::-webkit-scrollbar-corner {
  background-color: #646464;
}
::-webkit-resizer {
  background-color: #666;
}

#app {
  font-family: "Montserrat", sans-serif;
  /* -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale; */
  font-size: 14px;
  margin: 0;
  padding: 0;
  color: $defaultTextColor;
  display: flex;
  flex-direction: column;
  height: 100%;

  .container-wrapper {
    flex-grow: 1;
    display: flex;

    .container-left {
      flex-basis: $sidebarWidth;
      min-width: $sidebarWidth;
      background-color: $bg-color-light;
      display: flex;
      flex-direction: column;
      padding: 4px;
      border-right: solid 4px $bg-color-dark;

      .container-minimap {
        flex-basis: 250px;
        margin-bottom: 2px;
      }

      .container-targets {
        flex-grow: 1;
        flex-basis: 100px; /* This can be any number; forces div to respect flex box height. */
        min-height: 100px;
        margin-top: 2px;
      }
    }

    .container-center {
      flex-grow: 1;
      display: flex;
      flex-direction: column;
      position: relative;
      padding: 4px 2px;

      .container-maintext {
        flex-grow: 1;
        margin-bottom: 2px;
      }

      .container-input {
        flex-shrink: 1;
        margin-top: 2px;
        margin-bottom: 2px;
      }

      .container-bars {
        margin-top: 2px;
      }
    }

    .container-right {
      flex-basis: $sidebarWidth;
      min-width: $sidebarWidth;
      background-color: $bg-color-light;
      display: flex;
      flex-direction: column;
      padding: 4px;
      border-left: solid 4px $bg-color-dark;

      .container-skills {
        flex-grow: 1;
        margin-bottom: 2px;
      }

      .container-inventory {
        flex-basis: 365px;
        margin-top: 2px;
      }
    }
  }

  .status-bar-container {
    flex-basis: 30px;
    position: relative;
    background-color: $bg-color;
    padding: 2px;
    margin-top: 2px;
    background-image: url(/gfx/status-bg-01.png);
  }
}
</style>
