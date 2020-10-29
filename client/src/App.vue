<template>
    <div id="app">
        <div class="container-wrapper">
            <div class="container-left">
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
            <div class="container-right">
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
    </div>
</template>

<script>
import { mapState } from 'vuex';
import MainText from '@/components/MainText';
import InputBox from '@/components/InputBox';
import Minimap from '@/components/Minimap';
import RoomTargets from '@/components/RoomTargets';
import Inventory from '@/components/Inventory';
import Vitals from '@/components/Vitals';
import Skills from '@/components/Skills';
import StatusBar from '@/components/StatusBar';
import ItemTooltip from '@/components/ItemTooltip';

export default {
    name: 'App',
    components: {
        InputBox,
        MainText,
        Minimap,
        RoomTargets,
        Inventory,
        Vitals,
        Skills,
        StatusBar,
        ItemTooltip
    },
    data: () => {
        return {
            windowHeight: 0,
            windowWidth: 0,
        }
    },
    computed: mapState(['allowGlobalHotkeys', 'objectEditorOpen', 'isConnected', 'playerInfo']),
    watch: {
        isConnected: function(connected) {
            let token = this.$store.state.autoLoginToken;
            if (connected) {
                this.$store.dispatch('showText', { data: `Welcome to Armeria!\n\n` });

                if (token.length > 0) {
                    const char = token.split(':')[0];
                    this.$store.dispatch('showText', { data: `You are automatically being logged in as '${char}'.\n` });
                    this.$store.dispatch('sendSlashCommand', {
                        command: `/logintoken ${token}`
                    });
                } else {
                    this.$store.dispatch('showText', { data: 'If you have an existing character, you can <b>/login</b>. Otherwise, <b>/create</b> a new one.\n' });
                }
            } else {
                window.document.title = '**DISCONNECTED** Armeria.io';
            }
        },

        playerInfo: function(info) {
            window.document.title = `${info.name} - Armeria.io`;
        }
    },
    methods: {
        onWindowResize() {
            this.windowHeight = window.innerHeight;
            this.windowWidth = window.innerWidth;

            document.querySelector('.container-center').style.maxWidth = `${this.windowWidth-500}px`;
        },

        onKeyUp(event) {
            if (!this.allowGlobalHotkeys) {
                return;
            }

            let moveCommand = '';

            switch(event.key.toLowerCase()) {
                case 'w':
                    moveCommand = "/move north";
                    break;
                case 'a':
                    moveCommand = "/move west";
                    break;
                case 's':
                    moveCommand = "/move south";
                    break;
                case 'd':
                    moveCommand = "/move east";
                    break;
                case 'q':
                    moveCommand = "/move down";
                    break;
                case 'e':
                    moveCommand = "/move up";
                    break;
                case 'escape':
                    this.$store.dispatch('setObjectTarget', '');
                    break;
                case 'enter':
                    this.$store.dispatch('setForceInputFocus', { forced: true });
                    break;
                case '/':
                    this.$store.dispatch('setForceInputFocus', { forced: true, text: '/' });
                    break;
            }

            if (moveCommand.length > 0) {
                if (event.shiftKey && this.$store.state.permissions.indexOf('CAN_BUILD')) {
                    moveCommand = moveCommand.replace('/move', '/room create');
                }

                this.$store.dispatch('sendSlashCommand', {
                    command: moveCommand
                });
            }
        }
    },

    mounted() {
        this.onWindowResize();

        window.addEventListener(
            'resize',
            this.onWindowResize
        );

        window.addEventListener(
            'keyup',
            this.onKeyUp
        );
    },

    destroyed() {
        window.removeEventListener(
            'resize',
            this.onWindowResize
        );
    }
}
</script>

<style lang="scss">
@import url('https://fonts.googleapis.com/css?family=Montserrat:100,100i,200,200i,300,300i,400,400i,500,500i,600,600i,700,700i,800,800i,900,900i');
@import url('https://fonts.googleapis.com/css?family=Inconsolata:400,700&display=swap');
$backgroundNormal: #111;
$backgroundLight: #1b1b1b;
$sidebarWidth: 250px;

html, body {
    padding: 0;
    margin: 0;
    height: 100%;
    background-color: $backgroundNormal;
    user-select: none;
}

::-webkit-scrollbar { width: 3px; height: 3px; }
::-webkit-scrollbar-button {  background-color: #666; }
::-webkit-scrollbar-track {  background-color: #646464; }
::-webkit-scrollbar-track-piece { background-color: #111; }
::-webkit-scrollbar-thumb { height: 50px; background-color: #333; border-radius: 0px; }
::-webkit-scrollbar-corner { background-color: #646464; }
::-webkit-resizer { background-color: #666; }

#app {
    font-family: 'Montserrat', sans-serif;
    /* -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale; */
    font-size: 14px;
    margin: 0;
    padding: 0;
    color: #dbe3e6;
    display: flex;
    flex-direction: column;
    height: 100%;

.container-wrapper {
    flex-grow: 1;
    display: flex;

.container-left {
    flex-basis: $sidebarWidth;
    min-width: $sidebarWidth;
    background-color: $backgroundLight;
    display: flex;
    flex-direction: column;

.container-minimap {
    flex-basis: 250px;
}

.container-targets {
    flex-grow: 1;
    min-height: 100px;
}
}

.container-center {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    box-shadow: 0px 0px 12px 1px #080808;
    position: relative;

.container-maintext {
    flex-grow: 1;
}

.container-input {
    flex-shrink: 1;
}

.container-bars {
    flex-basis: 30px;
}
}

.container-right {
    flex-basis: $sidebarWidth;
    min-width: $sidebarWidth;
    background-color: $backgroundLight;
    display: flex;
    flex-direction: column;

.container-skills {
    flex-grow: 1;
}

.container-inventory {
    flex-basis: 365px;
}
}
}

.status-bar-container {
    flex-basis: 30px;
    border-top: 1px solid #333;
    position: relative;
    background-color: #0e0e0e;
}
}
</style>
