<template>
  <div class="container" :class="{ active: isFocused }">
    <input
      class="input-box"
      ref="inputBox"
      type="text"
      v-model="textToSend"
      @keyup.enter="handleSendText"
      @keyup.escape="handleRemoveFocus"
      @focus="handleFocus"
      @blur="handleBlur"
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
      checkDebugCommands(cmd) {
        if (cmd === '//openeditor') {
          this.$store.dispatch('setObjectEditorOpen', true);
          return true;
        } else if (cmd === '//closeeditor') {
          this.$store.dispatch('setObjectEditorOpen', false);
          return true;
        }

        return false;
      },

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
        } else if (this.checkDebugCommands(slashCommand)) {
          // do nothing
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

<style scoped>
  .container {
    background-color: #0c0c0c;
  }

  .input-box {
    background-color: #0c0c0c;
    border: 0;
    height: 35px;
    width: 99%;
    color: #fff;
    font-family: 'Montserrat', sans-serif;
    font-weight: 500;
    font-size: 13px;
    padding-left: 5px;
  }

  .container.active {
    background-color: #000;
  }

  .container.active .input-box {
    background-color: #000;
  }

  .input-box:focus {
    outline: none;
  }
</style>