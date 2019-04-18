<template>
  <div class="container">
    <input
      class="input-box"
      type="text"
      v-model="textToSend"
      v-on:keyup.enter="handleSendText"
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
  methods: {
    handleSendText() {
      let slashCommand = this.textToSend;

      // If command doesn't start with /, assume using /say
      if (slashCommand.substr(0, 1) !== '/') {
        slashCommand = `/say ${slashCommand}`;
      }

      if (slashCommand.length > 0) {
        this.$store.dispatch('sendSlashCommand', {
          command: slashCommand
        });
      }

      this.textToSend = '';
    },

    handleFocus() {
      this.isFocused = true;
    },

    handleBlur() {
      this.isFocused = false;
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