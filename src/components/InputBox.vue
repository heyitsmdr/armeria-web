<template>
    <div class="container" :class="{ active: isFocused }">
        <input
                class="input-box"
                ref="inputBox"
                type="text"
                v-model="textToSend"
                @keyup.enter="handleSendText"
                @keyup.escape="handleRemoveFocus"
                @keyup="handleKeyUp"
                @keydown="handleKeyDown"
                @focus="handleFocus"
                @blur="handleBlur"
        />
        <div class="hotkey-overlay" v-if="!isFocused" @click="handleHotkeyOverlayClick">
            Hotkey Mode -- Press ENTER for Input Mode
        </div>
        <div class="command-helper-overlay" ref="commandHelper" v-if="commandHelpVisible">
            <div v-html="helpHTML"></div>
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
                commandHelpVisible: false,
                helpHTML: '',
            }
        },
        computed: mapState(['objectEditorOpen', 'forceInputFocus', 'commandHistory', 'commandDictionary']),
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

            // test "this is" yes
            getCommandSegments(str) {
                let args = [];
                let recording = '';
                let inQuotes = false;

                for(let i = 0; i < str.length; i++) {
                    const char = str[i];

                    if (char === ' ' && !inQuotes) {
                        if (recording.length > 0) {
                            args.push(recording);
                            recording = '';
                        }
                    } else if (char === '"') {
                        if (inQuotes) {
                            inQuotes = false;
                            args.push(recording);
                            recording = '';
                        } else {
                            inQuotes = true;
                        }
                    } else {
                        recording += char;
                    }
                }

                // If there is anything in the buffer, push it as an element.
                if (recording.length > 0) {
                    args.push(recording);
                }

                // If string ends in a space or open quotes, add an empty element to the array.
                if (!inQuotes && str.substr(str.length - 1, 1) === ' ') {
                    args.push('');
                }

                if (inQuotes && str.substr(str.length - 1, 1) === '"') {
                    args.push('');
                }

                return args;
            },

            renderHelp() {
                const rawCommand = this.textToSend.substr(1);
                const commandSegments = this.getCommandSegments(rawCommand);
                const baseCommand = commandSegments[0].toLowerCase();

                this.helpHTML = '';
                this.commandDictionary.forEach(cmd => {
                    if (baseCommand.length > cmd.name.length) {
                        return;
                    } else if (cmd.name.substr(0, baseCommand.length) !== baseCommand) {
                        return;
                    }

                    if (commandSegments.length > 1 && cmd.args && cmd.args.length > 0) {
                        // Arguments on a root-level command.
                        this.helpHTML += `<div><b><span style="color:#ffe500">/${cmd.name}</span></b> `;
                        let argHelp = '';
                        for(let i = 0; i < cmd.args.length; i++) {
                            const arg = cmd.args[i];
                            const bracketOpen = arg.Optional ? '<' : '[';
                            const bracketClose = arg.Optional ? '>' : ']';

                            if ((i + 1) <= (commandSegments.length - 1)) {
                                this.helpHTML += `<span style="color:#ffe500">${bracketOpen}${arg.Name}${bracketClose}</span> `;
                                argHelp = arg.Help;
                            } else {
                                this.helpHTML += `${bracketOpen}${arg.Name}${bracketClose} `;
                            }
                        }
                        if (argHelp.length > 0) {
                            this.helpHTML += ` - ${argHelp}`;
                        }
                        this.helpHTML += `</div>`;
                    } else if (commandSegments.length === 2 && cmd.subCommands && cmd.subCommands.length > 0) {
                        // Sub-commands.
                        this.helpHTML += `<div><span style="color:#ffe500"><b>/${cmd.name}</b> &lt;sub-command&gt;</span></div>`;
                        this.helpHTML += `<br><div><b>Sub-commands:</b></div>`;
                        for(let i = 0; i < cmd.subCommands.length; i++) {
                            const subcmd = cmd.subCommands[i];

                            if (commandSegments[1].length > subcmd.name.length) {
                                continue;
                            } else if (subcmd.name.substr(0, commandSegments[1].length) !== commandSegments[1]) {
                                continue;
                            }

                            this.helpHTML += `<div>&nbsp;&nbsp;<b><span style="color:#ffe500">${commandSegments[1]}</span>${subcmd.name.substr(commandSegments[1].length)}</b> - ${subcmd.help}</div>`;
                        }
                    } else if (commandSegments.length > 2 && cmd.subCommands && cmd.subCommands.length > 0) {
                        // Arguments on a sub-command.
                        this.helpHTML += `<div><b><span style="color:#ffe500">/${cmd.name}</span></b> `;
                        let subcmd = null
                        for(let i = 0; i < cmd.subCommands.length; i++) {
                            if (commandSegments[1] === cmd.subCommands[i].name) {
                                subcmd = cmd.subCommands[i];
                                break;
                            }
                        }
                        let argHelp = '';
                        if (subcmd) {
                            this.helpHTML += `<b><span style="color:#ffe500">${subcmd.name}</span></b> `;
                            if (subcmd.args && subcmd.args.length > 0) {
                                for(let i = 0; i < subcmd.args.length; i++) {
                                    const arg = subcmd.args[i];
                                    const bracketOpen = arg.Optional ? '<' : '[';
                                    const bracketClose = arg.Optional ? '>' : ']';

                                    if ((i + 1) <= (commandSegments.length - 2)) {
                                        this.helpHTML += `<span style="color:#ffe500">${bracketOpen}${arg.Name}${bracketClose}</span> `;
                                        argHelp = arg.Help;
                                    } else {
                                        this.helpHTML += `${bracketOpen}${arg.Name}${bracketClose} `;
                                    }
                                }
                            }
                        }
                        if (argHelp.length > 0) {
                            this.helpHTML += ` - ${argHelp}`;
                        }
                        this.helpHTML += `</div>`;
                    } else {
                        this.helpHTML += `<div>` +
                            `<b><span style="color:#ffe500">/${baseCommand}</span>${cmd.name.substr(baseCommand.length)}</b>` +
                            ` - ${cmd.help}` +
                            `</div>`;
                    }
                });

                // Show or hide depending on results being found.
                if (this.helpHTML.length === 0) {
                    this.commandHelpVisible = false;
                } else {
                    this.commandHelpVisible = true;
                    this.$nextTick(() => {
                        const commandHelperHeight = this.$refs['commandHelper'].clientHeight;
                        this.$refs['commandHelper'].style.top = `-${commandHelperHeight + 2}px`;
                    });
                }
            },

            handleSendText() {
                let slashCommand = this.textToSend;

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
                this.commandHelpVisible = false;
                this.$nextTick(() => {
                    event.target.blur();
                });
            },

            handleFocus() {
                this.isFocused = true;
                this.$store.dispatch('setAllowGlobalHotkeys', false);
            },

            handleBlur() {
                this.$store.dispatch('setAllowGlobalHotkeys', true);
                this.commandHelpVisible = false;
                this.$nextTick(() => {
                    this.isFocused = false;
                });
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

            handleKeyUp() {
                // Render help for slash commands.
                if (this.textToSend.substr(0, 1) === '/' && this.textToSend.length > 1) {
                    this.renderHelp();
                } else {
                    this.commandHelpVisible = false;
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
        background: #000;
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

    .command-helper-overlay {
        position: absolute;
        background: rgb(0,0,0);
        background: linear-gradient(180deg, rgba(0,0,0,0) 0%, rgba(0,0,0,0.9) 60%);
        width: 99%;
        padding: 20px 5px 10px 5px;
        font-size: 12px;
    }
</style>