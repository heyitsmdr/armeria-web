<template>
    <div class="container" :class="{ active: isFocused }">
        <input
                class="input-box"
                ref="inputBox"
                type="text"
                v-model="textToSend"
                @keyup.enter="handleSendText"
                @keyup.escape="handleRemoveFocus"
                @keydown="handleKeyDown"
                @focus="handleFocus"
                @blur="handleBlur"
        />
        <div class="hotkey-overlay" v-if="!isFocused" @click="handleHotkeyOverlayClick">
            Hotkey Mode -- Press ENTER for Input Mode
        </div>
    </div>
</template>

<script>
    import {mapState} from 'vuex';

    export default {
        name: 'InputBox',
        data: () => {
            return {
                textToSend: '',
                password: '',
                isFocused: false,
                lastCommandHistoryIndex: -1,
            }
        },
        computed: mapState(['objectEditorOpen', 'forceInputFocus', 'commandHistory']),
        mounted() {
            this.$refs['inputBox'].focus();
        },
        watch: {
            forceInputFocus: function (data) {
                if (data.forced) {
                    this.$refs['inputBox'].focus();
                    if (data.text) {
                        this.textToSend = data.text;
                    }
                    this.$store.dispatch('setForceInputFocus', {forced: false});
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
                } else if (cmd === '//clearttcache') {
                    this.$store.dispatch('clearItemTooltipCache', false);
                    this.$store.dispatch('showText', { data: `\n[DEBUG] Item tooltip cache has been cleared on your client.\n` });
                    return true;
                }

                return false;
            },

            selectAll: function() {
                this.$refs['inputBox'].select();
            },

            getLastCommand() {
                let retrieveIndex = 0;

                if (this.lastCommandHistoryIndex === -1) {
                    retrieveIndex = this.commandHistory.length - 1;
                    this.lastCommandHistoryIndex = retrieveIndex;
                } else if (this.lastCommandHistoryIndex > 0) {
                    retrieveIndex = this.lastCommandHistoryIndex - 1;
                    this.lastCommandHistoryIndex = retrieveIndex
                }

                return this.commandHistory[retrieveIndex];
            },

            getNextCommand() {
                let retrieveIndex = this.lastCommandHistoryIndex;

                if (retrieveIndex === -1) {
                    retrieveIndex = this.commandHistory.length - 1;
                    this.lastCommandHistoryIndex = retrieveIndex;
                } else if (this.lastCommandHistoryIndex < (this.commandHistory.length - 1)) {
                    retrieveIndex = this.lastCommandHistoryIndex + 1;
                    this.lastCommandHistoryIndex = retrieveIndex
                }

                return this.commandHistory[retrieveIndex];
            },

            handleSendText() {
                let slashCommand = this.textToSend;

                this.$store.dispatch('showText', {
                    data: `<div class="inline-loopback">${slashCommand}</div>`
                });

                if (slashCommand.length === 0) {
                    this.$store.dispatch('sendSlashCommand', {
                        command: '/look'
                    });
                } else if (slashCommand.substr(0, 1) !== '/') {
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
                this.lastCommandHistoryIndex = -1;
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
            },

            handleKeyDown(e) {
                if (e.key === 'ArrowUp') {
                    this.textToSend = this.getLastCommand();
                    setTimeout(this.selectAll, 10);
                } else if (e.key === 'ArrowDown') {
                    this.textToSend = this.getNextCommand();
                    setTimeout(this.selectAll, 10);
                } else if (this.textToSend.substr(0, 6).toLowerCase() === '/login' && this.textToSend.split(" ").length === 3) {
                    if (e.key === 'Backspace') {
                        this.password = this.password.slice(0, this.password.length - 1);
                        this.textToSend = this.textToSend.slice(0, this.textToSend.length - 1);
                        this.textToSend += "*";
                    } else if (e.key !== 'Enter' && e.key !== 'Escape') {
                        e.preventDefault();
                        e.stopPropagation();
                        this.password += e.key;
                        this.textToSend += "*";
                    }
                } else {
                    this.password = "";
                }
            },

            handleHotkeyOverlayClick() {
                this.$refs['inputBox'].focus();
            }
        }
    }
</script>

<style scoped>
    .container {
        border: 1px solid #222;
        position: relative;
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
        border: 1px solid #adadad;
    }

    .container.active .input-box {
        background-color: #000;
    }

    .input-box:focus {
        outline: none;
    }

    .hotkey-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        z-index: 10;
        height: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        background-color: rgba(0, 0, 0, 0.8);
        color: #888;
    }

    .hotkey-overlay:hover {
        cursor: pointer;
        color: #bbb;
    }
</style>