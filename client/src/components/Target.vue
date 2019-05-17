<template>
    <div>
        <div
                class="container"
                :class="{'selected': this.objectTarget == this.name}"
                ref="container"
                @mousedown="onMouseDown"
                @mouseup="onMouseUp"
                @contextmenu.stop.prevent="onContextMenu"
        >
            <div class="picture">
                <div class="picture-container"
                     @dragenter.stop.prevent="onPictureDragEnter"
                     @drop.stop.prevent="onPictureDragDrop"
                     @dragleave.stop.prevent
                     @dragover.stop.prevent
                >

                </div>
            </div>
            <div class="name">
                <div class="name-container">
                    {{ name }}
                </div>
            </div>
        </div>
    </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
    name: 'Target',
    props: ['name'],
    methods: {
        onMouseDown: function() {
            this.$refs['container'].classList.add('mouse-down');
        },

        onMouseUp: function() {
            this.$refs['container'].classList.remove('mouse-down');
            this.$store.dispatch('setObjectTarget', this.name);
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
    },
    computed: mapState(['objectTarget'])
}
</script>

<style lang="scss" scoped>
.container {
    background-color: #0c0c0c;
    border: 1px solid #353535;
    margin: 0 10px 10px 10px;
    transition: all .1s ease-in-out;
    transform: scale(1);
    display: flex;

    &.selected {
         border: 1px solid #ffeb3b;
    }

    &.mouse-down {
        transform: scale(1.01) !important;
    }

    &:hover {
        cursor: pointer;
        transform: scale(1.05);
    }

    .picture {
        flex-basis: 50px;

        .picture-container {
             height: 50px;
             box-shadow: inset 0px 0px 15px #000;
        }
    }

    .name {
        flex-grow: 1;
        display: flex;
        align-items: center;
        margin-left: 10px;

        .name-container {
            font-weight: 600;
        }
    }
}
</style>