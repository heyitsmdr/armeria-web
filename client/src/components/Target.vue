<template>
    <div>
        <div
                class="container"
                :class="{'selected': objectTarget === name, 'is-character': objectType === 0,  'is-mob': objectType == 1, 'is-item': objectType == 2 }"
                ref="container"
                @mousedown="handleMouseDown"
                @mouseup="handleMouseUp"
                @dblclick="handleDoubleClick"
                @contextmenu.stop.prevent="onContextMenu"
        >
            <div class="picture">
                <div class="picture-container"
                     :style="{ backgroundImage: getBackgroundUrl() }"
                     @dragenter.stop.prevent="onPictureDragEnter"
                     @drop.stop.prevent="onPictureDragDrop"
                     @dragleave.stop.prevent
                     @dragover.stop.prevent
                >

                </div>
            </div>
            <div class="name">
                <div class="name-container">
                    <div>{{ name }}</div>
                    <div class="alt">{{ title }}</div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
    name: 'Target',
    props: ['name', 'objectType', 'pictureKey', 'title'],
    computed: mapState(['isProduction', 'objectTarget']),
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
            if (e.shiftKey) {
                if (this.objectType === 0) {
                    this.$socket.sendObj({ type: 'command', payload: '/character edit ' + this.name });
                } else if (this.objectType === 1) {
                    this.$socket.sendObj({ type: 'command',  payload: '/mob edit ' + this.name });
                } else if (this.objectType === 2) {
                    this.$socket.sendObj({ type: 'command',  payload: '/item edit ' + this.name });
                }
            } else {
                this.$store.dispatch('setObjectTarget', this.name);
            }
        },

        handleDoubleClick: function() {
            if (this.objectType === 2) {
                this.$socket.sendObj({ type: 'command', payload: '/get ' + this.name });
                this.$store.dispatch('setObjectTarget', '');
            }
        },

        onPictureDragEnter: function(event) {
            // TODO: add class to make it obvious you can drop something here
        },

        onPictureDragDrop: function(event) {
            const files = event.dataTransfer.files;

            if (files.length > 1) {
                this.$store.dispatch('showText', { data: '\nYou can only upload one image at a time.' });
                return;
            }

            const file = files[0];

            console.log(file);

            console.log('drag drop',event.dataTransfer.files);
            [...event.dataTransfer.files].forEach(f => {
                const reader = new FileReader()
                reader.readAsBinaryString(f)
                reader.onload = function() {
                    console.log(btoa(reader.result));
                }
            });
        },

        onContextMenu: function(event) {
            // TODO: Add a custom right-click menu
            // https://dev.to/iamafro/how-to-create-a-custom-context-menu--5d7p
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

    &.selected {
         border: 1px solid #ffeb3b !important;
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
    }
}
</style>