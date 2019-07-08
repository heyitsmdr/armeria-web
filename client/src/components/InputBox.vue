<template>
  <div class="container" :class="{ active: isFocused }">
    <input
      class="input-box"
      ref="inputBox"
      type="text"
      v-model="textToSend"
      @keyup.enter="handleSendText"
      @keyup.escape="handleRemoveFocus"
      @keypress="handleKeyPress"
      @focus="handleFocus"
      @blur="handleBlur"
    />
  </div>
</template>

<script>
  import { mapState} from 'vuex';

  export default {
    name: 'InputBox',
    data: () => {
      return {
        textToSend: '',
        password: '',
        isFocused: false,
      }
    },
    computed: mapState(['objectEditorOpen', 'forceInputFocus']),
    mounted() {
      this.$refs['inputBox'].focus();
    },
    watch: {
        forceInputFocus: function(force) {
            if (force) {
                this.$refs['inputBox'].focus();
                this.$store.dispatch('setForceInputFocus', false);
            }
        }
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
        } else if (slashCommand.substr(0, 6).toLowerCase() === '/login') {
          let characterName = slashCommand.split(' ')[1];
          this.$store.dispatch('sendSlashCommand', {
            command: `/login ${characterName} ${this.password}`
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
        if (this.objectEditorOpen) {
          this.$store.dispatch('setObjectEditorOpen', false);
          return;
        }

        event.target.blur()
      },

      handleFocus() {
        this.isFocused = true;
        this.$store.dispatch('setAllowGlobalHotkeys', false);
      },

      handleBlur() {
        this.isFocused = false;
        this.$store.dispatch('setAllowGlobalHotkeys', true);
      },

      // TODO: Add functionality to be able to handle a backspace press (might only be on keyDown)
      handleKeyPress(e) {
        if (this.textToSend.substr(0, 6).toLowerCase() === '/login' && this.textToSend.split(" ").length === 3) {
          if (e.key !== 'Enter') {
            e.preventDefault();
            e.stopPropagation();
            this.password += e.key;
            this.textToSend += "*";
          }
        } else {
          this.password = "";
        }
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
    border: 1px solid #303030;
    box-shadow: 0px 0px 10px 0px #292929;
  }

  .container.active .input-box {
    background-color: #000;
  }

  .input-box:focus {
    outline: none;
  }
</style>