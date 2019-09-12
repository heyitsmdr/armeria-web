<template>
    <div>
        <div
                class="container"
                ref="container"
        >
            <div class="picture">
                <div class="picture-container"
                     :style="{ backgroundImage: getBackgroundUrl() }"
                >

                </div>
            </div>
            <div class="name">
                <div class="name-container">
                    <div>{{ name }}</div>
                    <div class="alt">{{ title }}</div>
                </div>
                <div class="you" :class="{ selected: uuid===objectTargetUUID }" v-if="uuid===playerInfo.uuid">you</div>
            </div>
            <div
                class="overlay"
                @mousedown="handleMouseDown"
                @mouseup="handleMouseUp"
                @dblclick="handleDoubleClick"
                @contextmenu.stop.prevent="onContextMenu"
                @dragenter="handleDragEnter"
                @dragleave="handleDragLeave"
                @drop="handleDrop"
                @dragover.prevent
            ></div>
        </div>
    </div>
</template>

<script>
const OBJECT_TYPE_CHARACTER = 0;
const OBJECT_TYPE_MOB = 1;
const OBJECT_TYPE_ITEM = 2;
import { mapState } from 'vuex';
import {PICKUP_ITEM} from "../plugins/SFX";

export default {
    name: 'Target',
    props: ['uuid', 'name', 'objectType', 'pictureKey', 'title', 'color'],
    computed: mapState(['isProduction', 'objectTargetUUID', 'playerInfo']),
    watch: {
        objectTargetUUID: function(target) {
            if (this.uuid === target) {
                this.$refs['container'].classList.add('selected');
            } else {
                this.$refs['container'].classList.remove('selected');
            }
        },
    },
    mounted() {
        switch(this.objectType) {
            case OBJECT_TYPE_CHARACTER:
                this.$refs['container'].classList.add('is-character');
                break;
            case OBJECT_TYPE_MOB:
                this.$refs['container'].classList.add('is-mob');
                break;
            case OBJECT_TYPE_ITEM:
                this.$refs['container'].classList.add('is-item');
                if (this.color.length > 0) {
                    this.$refs['container'].style.borderColor = this.color;
                }
                break;
        }

        if (this.uuid === this.objectTargetUUID) {
            this.$refs['container'].classList.add('selected');
        }
    },
    methods: {
        getBackgroundUrl() {
            if (!this.isProduction) {
                return `url(http://localhost:8081/oi/${this.pictureKey})`;
            }

            return `url(/oi/${this.pictureKey})`;
        },

        handleMouseDown: function() {
            this.$refs['container'].classList.add('mouse-down');
        },

        handleMouseUp: function(e) {
            this.$refs['container'].classList.remove('mouse-down');
            if (this.$store.state.permissions.indexOf('CAN_BUILD') >= 0) {
                if (e.shiftKey) {
                    if (this.objectType === OBJECT_TYPE_CHARACTER) {
                        this.$socket.sendObj({ type: 'command', payload: '/character edit ' + this.name });
                    } else if (this.objectType === OBJECT_TYPE_MOB) {
                        this.$socket.sendObj({ type: 'command',  payload: '/mob iedit ' + this.uuid });
                    } else if (this.objectType === OBJECT_TYPE_ITEM) {
                        this.$socket.sendObj({ type: 'command',  payload: '/item iedit ' + this.uuid });
                    }
                    return
                }
            }

            this.$store.dispatch('setObjectTarget', this.uuid);
        },

        handleDoubleClick: function() {
            if (this.objectType === OBJECT_TYPE_ITEM) {
                this.$socket.sendObj({ type: 'command', payload: '/get ' + this.name });
                this.$store.dispatch('setObjectTarget', '');
                this.$soundEvent(PICKUP_ITEM);
            }
        },

        onContextMenu: function() {
            this.$socket.sendObj({ type: 'command', payload: '/look ' + this.uuid });
        },

        handleDragEnter: function() {
            this.$refs['container'].classList.add('can-drop-item');
        },

        handleDragLeave: function() {
            this.$refs['container'].classList.remove('can-drop-item');
        },

        handleDrop: function(e) {
            this.$refs['container'].classList.remove('can-drop-item');
            const item_uuid = e.dataTransfer.getData('item_uuid');
            if (item_uuid) {
                this.$store.dispatch('sendSlashCommand', {
                    command: `/give ${this.uuid} ${item_uuid}`
                });
            }
        }
    }
}
</script>

<style lang="scss" scoped>
.container {
    background-color: #0c0c0c;
    margin: 0 10px 10px 10px;
    transition: all .1s ease-in-out;
    transform: scale(1);
    display: flex;

    &.can-drop-item {
         transform: scale(1.1) !important;
    }

    &.selected {
         border: 1px solid #ffeb3b !important;
         background-color: #231f00;
    }

    &.mouse-down {
        transform: scale(1.01) !important;
    }

    &.is-character {
         border: 1px solid #353535;
    }

    &.is-mob {
        border: 1px solid #673604;

        .name {
            color: #d48a3e;
        }
    }

    &.is-item {
         border: 1px solid #fff;

        .name {
            color: #fff;
        }
    }

    &:hover {
        cursor: pointer;
        transform: scale(1.05);
    }

    .picture {
        flex-basis: 50px;

        .picture-container {
            height: 50px;
            box-shadow: inset 0px 0px 5px 0px #3a3a3a;
            background-size: contain;
        }
    }

    .name {
        flex-grow: 1;
        display: flex;
        align-items: center;
        margin-left: 10px;

        .name-container {
            font-weight: 600;

            .alt {
                font-weight: 400;
                font-size: 12px;
            }
        }

        .you {
            position: absolute;
            right: -1px;
            top: -1px;
            background-color: #353535;
            padding: 2px 5px;
            border: 1px solid #353535;
            text-transform: uppercase;
            font-size: 12px;
            transition: all .1s ease-in-out;

            &.selected {
                 background-color: #eedb38;
                 border: 1px solid #eedb38;
                 color: #000;
            }
        }
    }

    .overlay {
        position: absolute;
        top: 0px;
        left: 0px;
        height: 100%;
        width: 100%;
        z-index: 999;
    }
}
</style>