<template>
  <div class="root">
    <div class="inner-container">
      <div class="ping">Ping: {{ pingTime }}ms</div>
      <div class="version">
        {{ deployAppName }}
        <span v-if="normalizedDeployVersion.length > 0"
          >({{ normalizedDeployVersion }})</span
        >
      </div>
    </div>
  </div>
</template>

<script>
import { mapState, mapGetters } from "vuex";

export default {
  name: "StatusBar",
  computed: {
    ...mapState(["pingTime", "deployAppName"]),
    ...mapGetters(["normalizedDeployVersion"]),
  },
};
</script>

<style scoped lang="scss">
@import "@/styles/common";

.root {
  display: flex;
  height: 100%;
  align-items: center;
  font-size: 12px;
  box-sizing: border-box;

  .inner-container {
    display: flex;
    width: 100%;
    padding: 0 5px;

    .ping {
      flex-shrink: 1;
    }

    .version {
      flex-grow: 1;
      text-align: right;
    }
  }
}

.ping {
}
</style>
