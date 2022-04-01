<template>
    <div class="root" :class="{ active: isFocused }">
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
            Hotkey Mode Active. Click/Enter to type.
        </div>
        <div class="command-helper-overlay" ref="commandHelper" v-if="commandHelpVisible">
            <div v-html="helpHTML"></div>
        </div>
    </div>
</template>

<script setup>
    import { ref, computed, watch, nextTick, onMounted } from 'vue';
    import { useStore } from 'vuex';

    // Defaults.
    const inputBox = ref(null);      // Vue auto maps this to the underlying HTML reference.
    const commandHelper = ref(null); // Auto-mapped.
    const textToSend = ref('');
    const password = ref('');
    const isFocused = ref(false);
    const lastCommandHistoryIndex = ref(-1);
    const commandHelpVisible = ref(false);
    const helpHTML = ref('');

    // State from store.
    const store = useStore();
    const forceInputFocus = computed(() => store.state.forceInputFocus);
    const commandHistory = computed(() => store.state.commandHistory);
    const commandDictionary = computed(() => store.state.commandDictionary);

    // Computed.
    const expandedCommandDictionary = computed(() => {
        const dict = [];
        commandDictionary.value.forEach(d => dict.push(d));
        commandDictionary.value.forEach(cmd => {
            if (cmd.altNames) {
                cmd.altNames.forEach(alt => {
                    let newCmd = Object.assign({}, cmd);
                    newCmd.name = alt;
                    newCmd.altNames = [];
                    dict.push(newCmd);
                });

            }
        });
        return dict;
    });

    // Watches.
    watch(forceInputFocus, (newValue) => {
        if (newValue.forced) {
            inputBox.value.focus();
            if (newValue.text) {
                textToSend.value = newValue.text;
            }
            store.dispatch('setForceInputFocus', { forced: false });
        }
    });

    // Mounted.
    onMounted(() => {
        inputBox.value.focus();
    });

    /**
     * Checks if the player is using a debug command.
     * @param {String} cmd
     * @returns {Boolean}
     */
    function checkDebugCommands(cmd) {
        if (cmd === '//openeditor') {
            store.dispatch('setObjectEditorOpen', true);
            return true;
        } else if (cmd === '//closeeditor') {
            store.dispatch('setObjectEditorOpen', false);
            return true;
        } else if (cmd === '//clearttcache') {
            store.dispatch('clearItemTooltipCache', false);
            store.dispatch('showText', { data: `\n[DEBUG] Item tooltip cache has been cleared on your client.\n` });
            return true;
        }

        return false;
    }

    /**
     * Selects all the text in the input box.
     */
    function selectAll() {
        inputBox.value.select();
    }

    /**
     * Returns the last command used, when traversing the history.
     * @returns {String}
     */
    function getLastCommand() {
        let retrieveIndex = 0;

        if (lastCommandHistoryIndex.value === -1) {
            retrieveIndex = commandHistory.value.length - 1;
            lastCommandHistoryIndex.value = retrieveIndex;
        } else if (lastCommandHistoryIndex.value > 0) {
            retrieveIndex = lastCommandHistoryIndex.value - 1;
            lastCommandHistoryIndex.value = retrieveIndex
        }

        return commandHistory.value[retrieveIndex];
    }

    /**
     * Returns the next command used, when traversing the history.
     * @returns {*}
     */
    function getNextCommand() {
        let retrieveIndex = lastCommandHistoryIndex.value;

        if (retrieveIndex === -1) {
            retrieveIndex = commandHistory.value.length - 1;
            lastCommandHistoryIndex.value = retrieveIndex;
        } else if (lastCommandHistoryIndex.value < (commandHistory.value.length - 1)) {
            retrieveIndex = lastCommandHistoryIndex.value + 1;
            lastCommandHistoryIndex.value = retrieveIndex
        }

        return commandHistory.value[retrieveIndex];
    }

    /**
     * Returns the command broken down into segments.
     * @param {String} str
     * @returns {Array}
     */
    function getCommandSegments(str) {
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
    }

    /**
     * Renders the help pop-up, used while typing commands.
     * @returns {Promise<void>}
     */
    async function renderHelp() {
        const rawCommand = textToSend.value.substring(1);
        const commandSegments = getCommandSegments(rawCommand);
        const baseCommand = commandSegments[0].toLowerCase();

        helpHTML.value = '';
        for(let i = 0; i < expandedCommandDictionary.value.length; i++) {
            const cmd = expandedCommandDictionary.value[i];
            if (baseCommand.length > cmd.name.length) {
                continue;
            } else if (cmd.name.substring(0, baseCommand.length) !== baseCommand) {
                continue;
            }

            if (commandSegments.length > 1 && cmd.args && cmd.args.length > 0) {
                // Arguments on a root-level command.
                if (baseCommand !== cmd.name) {
                    continue;
                }
                helpHTML.value += `<div><b><span style="color:#ffe500">/${cmd.name}</span></b> `;
                let argHelp = '';
                for(let i = 0; i < cmd.args.length; i++) {
                    const arg = cmd.args[i];
                    const bracketOpen = arg.Optional ? '&lt;' : '[';
                    const bracketClose = arg.Optional ? '&gt;' : ']';

                    if ((i + 1) <= (commandSegments.length - 1)) {
                        helpHTML.value += `<span style="color:#ffe500">${bracketOpen}${arg.Name}${bracketClose}</span> `;
                        argHelp = arg.Help;
                    } else {
                        helpHTML.value += `${bracketOpen}${arg.Name}${bracketClose} `;
                    }
                }
                if (argHelp.length > 0) {
                    helpHTML.value += ` - ${argHelp}`;
                }
                helpHTML.value += `</div>`;
            } else if (commandSegments.length === 2 && cmd.subCommands && cmd.subCommands.length > 0) {
                // Sub-commands.
                helpHTML.value += `<div><span style="color:#ffe500"><b>/${cmd.name}</b> &lt;sub-command&gt;</span></div>`;
                helpHTML.value += `<br><div><b>Sub-commands:</b></div>`;
                for(let i = 0; i < cmd.subCommands.length; i++) {
                    const subcmd = cmd.subCommands[i];

                    if (commandSegments[1].length > subcmd.name.length) {
                        continue;
                    } else if (subcmd.name.substring(0, commandSegments[1].length) !== commandSegments[1]) {
                        continue;
                    }

                    helpHTML.value += `<div>&nbsp;&nbsp;<b><span style="color:#ffe500">${commandSegments[1]}</span>${subcmd.name.substr(commandSegments[1].length)}</b> - ${subcmd.help}</div>`;
                }
            } else if (commandSegments.length > 2 && cmd.subCommands && cmd.subCommands.length > 0) {
                // Arguments on a sub-command.
                helpHTML.value += `<div><b><span style="color:#ffe500">/${cmd.name}</span></b> `;
                let subcmd = null
                for(let i = 0; i < cmd.subCommands.length; i++) {
                    if (commandSegments[1] === cmd.subCommands[i].name) {
                        subcmd = cmd.subCommands[i];
                        break;
                    }
                }
                let argHelp = '';
                if (subcmd) {
                    helpHTML.value += `<b><span style="color:#ffe500">${subcmd.name}</span></b> `;
                    if (subcmd.args && subcmd.args.length > 0) {
                        for(let i = 0; i < subcmd.args.length; i++) {
                            const arg = subcmd.args[i];
                            const bracketOpen = arg.Optional ? '&lt;' : '[';
                            const bracketClose = arg.Optional ? '&gt;' : ']';

                            if ((i + 1) <= (commandSegments.length - 2)) {
                                helpHTML.value += `<span style="color:#ffe500">${bracketOpen}${arg.Name}${bracketClose}</span> `;
                                argHelp = arg.Help;
                            } else {
                                helpHTML.value += `${bracketOpen}${arg.Name}${bracketClose} `;
                            }
                        }
                    }
                }
                if (argHelp.length > 0) {
                    helpHTML.value += ` - ${argHelp}`;
                }
                helpHTML.value += `</div>`;
            } else {
                helpHTML.value += `<div>` +
                    `<b><span style="color:#ffe500">/${baseCommand}</span>${cmd.name.substring(baseCommand.length)}</b>` +
                    ` - ${cmd.help}` +
                    ` <span style="color:#f00;font-weight:600">${(cmd.permissions && cmd.permissions.RequirePermission) ? '['+cmd.permissions.RequirePermission+']' : ''}</span>` +
                    `</div>`;
            }
        }

        // Show or hide depending on results being found.
        if (helpHTML.value.length === 0) {
            commandHelpVisible.value = false;
        } else {
            commandHelpVisible.value = true;
            await nextTick();
            const commandHelperHeight = commandHelper.value.clientHeight;
            commandHelper.value.style.top = `-${commandHelperHeight + 2}px`;
        }
    }

    /**
     * Handles sending the text/command to the server.
     */
    function handleSendText() {
        let slashCommand = textToSend.value;

        if (slashCommand.length === 0) {
            store.dispatch('sendSlashCommand', {
                command: '/look',
                hidden: true,
            });
        } else if (slashCommand.substring(0, 1) !== '/') {
            store.dispatch('sendSlashCommand', {
                command: `/say ${slashCommand}`,
                hidden: true,
            });
        } else if (slashCommand.substring(0, 6).toLowerCase() === '/login') {
            let characterName = slashCommand.split(' ')[1];
            store.dispatch('sendSlashCommand', {
                command: `/login ${characterName} ${password.value}`
            });
        } else if (checkDebugCommands(slashCommand)) {
            // do nothing
        } else {
            store.dispatch('sendSlashCommand', {
                command: slashCommand
            });
        }


        textToSend.value = '';
        lastCommandHistoryIndex.value = -1;
    }

    /**
     * Forces the input box to lose focus (i.e. the player pressed ESC).
     * @param {KeyboardEvent} event
     * @returns {Promise<void>}
     */
    async function handleRemoveFocus(event) {
        commandHelpVisible.value = false;
        await nextTick();
        event.target.blur();
    }

    /**
     * Handles when the input box gains focus.
     */
    function handleFocus() {
        isFocused.value = true;
        store.dispatch('setAllowGlobalHotkeys', false);
    }

    /**
     * Handles when the input box loses focus.
     */
    async function handleBlur() {
        await store.dispatch('setAllowGlobalHotkeys', true);
        commandHelpVisible.value = false;
        await nextTick();
        isFocused.value = false;
    }

    /**
     * Handles when a key is pressed down.
     * @param {KeyboardEvent} e
     */
    function handleKeyDown(e) {
        if (e.key === 'ArrowUp') {
            textToSend.value = getLastCommand();
            setTimeout(selectAll, 10);
        } else if (e.key === 'ArrowDown') {
            textToSend.value = getNextCommand();
            setTimeout(selectAll, 10);
        } else if (textToSend.value.substr(0, 6).toLowerCase() === '/login' && textToSend.value.split(" ").length === 3) {
            if (e.key === 'Backspace') {
                password.value = password.value.slice(0, password.value.length - 1);
                textToSend.value = textToSend.value.slice(0, textToSend.value.length - 1);
                textToSend.value += "*";
            } else if (e.key !== 'Enter' && e.key !== 'Escape') {
                e.preventDefault();
                e.stopPropagation();
                password.value += e.key;
                textToSend.value += "*";
            }
        } else {
            password.value = "";
        }
    }

    /**
     * Handles when a key is depressed.
     */
    function handleKeyUp() {
        // Render help for slash commands.
        if (textToSend.value.substr(0, 1) === '/' && textToSend.value.length > 1) {
            renderHelp();
        } else {
            commandHelpVisible.value = false;
        }
    }

    /**
     * Handles when the hotkey overlay is clicked on.
     */
    function handleHotkeyOverlayClick() {
        inputBox.value.focus();
    }
</script>

<style scoped lang="scss">
    @import "@/styles/common";
    $height: 35px;
    
    .root {
        position: relative;
        background: $bg-color;
        padding-left: 5px;
        border: $defaultBorder;
    }

    .input-box {
        display:block;
        width:100%;
        padding:0;
        border-width:0;
        border: 0;
        height: $height;
        color: $defaultTextColor;
        font-family: 'Montserrat', sans-serif;
        font-weight: 500;
        font-size: 13px;
    }

    .root.active {
        border: $defaultBorder;
    }

    .root.active .input-box {
        background-color: $bg-color;
    }

    .input-box:focus {
        margin: 0;
        padding: 0;
        border: 0;
        outline: none;
    }

    .hotkey-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        z-index: 10;
        height: $height;
        display: flex;
        justify-content: center;
        align-items: center;
        background-color: $bg-color;
        color: $defaultTextColor;
    }

    .hotkey-overlay:hover {
        cursor: pointer;
        color: $defaultTextColor;
    }

    .command-helper-overlay {
        position: absolute;
        background: $bg-color;
        background: $bg-color;
        width: 99%;
        padding: 20px 5px 10px 5px;
        font-size: 12px;
    }
</style>