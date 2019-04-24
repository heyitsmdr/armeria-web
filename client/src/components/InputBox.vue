<template>
  <div class="container">
    <input
      class="input-box"
      ref="inputBox"
      type="text"
      v-model="textToSend"
      v-on:keyup.enter="handleSendText"
      v-on:keyup.escape="handleRemoveFocus"
      @focus="handleFocus"
      @blur="handleBlur"
      v-bind:class="{ active: isFocused }"
    />
  </div>
</template>

<script>
export default {
  name: 'InputBox',
  data: () => {
    return {
      textToSend: '',
      isFocused: false,
    }
  },
  mounted() {
    this.$refs['inputBox'].focus();
  },
  methods: {
    handleSendText() {
      let slashCommand = this.textToSend;

      if (slashCommand.length === 0) {
        this.$store.dispatch('sendSlashCommand', {
          command: '/look'
        });
      }
      else if (slashCommand.substr(0, 1) !== '/') {
        this.$store.dispatch('sendSlashCommand', {
          command: `/say ${slashCommand}`
        });
      } else {
        this.$store.dispatch('sendSlashCommand', {
          command: slashCommand
        });
      }

      this.textToSend = '';
    },

    handleRemoveFocus(event) {
      event.target.blur()
    },

    handleFocus() {
      this.isFocused = true;
      this.$store.dispatch('setAllowGlobalHotkeys', false);
    },

    handleBlur() {
      this.isFocused = false;
      this.$store.dispatch('setAllowGlobalHotkeys', true);
    }
  }
}
</script>

<style lang="scss" scoped>
.container {

}

.input-box {
  background-color: #0c0c0c;
  border: 0;
  height: 35px;
  width: 100%;
  color: #fff;
  font-family: 'Montserrat', sans-serif;
  font-weight: 500;
  font-size: 13px;
  padding-left: 5px;

  &.active {
    background-color: #000;
  }

  &:focus{
    outline: none;
  }
}
</style>