<template>
  <div class="container">
    <input class="input-box" type="text" v-model="textToSend" v-on:keyup.enter="handleSendText" />
  </div>
</template>

<script>
export default {
  name: 'InputBox',
  data: () => {
    return {
      textToSend: ''
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
    }
  }
}
</script>

<style lang="scss" scoped>
.container {
  padding-left: 5px;
}

.input-box {
  background-color: #000;
  border: 0;
  height: 35px;
  width: 100%;
  color: #fff;
  font-family: 'Montserrat', sans-serif;
  font-weight: 500;
  font-size: 13px;

  &:focus{
    outline: none;
  }
}
</style>