<template>
    <div class="root" @dragenter="handleDragEnter" @dragleave="handleDragLeave" @dragover.prevent @drop="handleDrop">
        <div class="drop-overlay" ref="overlay">
            Drop the item here to place the item into the room.
        </div>
        <div class="container" @click="handleClick">
            <Target
                v-for="obj in sortedRoomObjects"
                :key="obj.uuid"
                :name="obj.name"
                :pictureKey="obj.picture"
                :objectType="obj.type"
                :title="obj.title"
                :color="obj.color"
            />
        </div>
    </div>
</template>

<script>
import { mapState } from 'vuex';
import Target from '@/components/Target';

export default {
    name: 'RoomTargets',
    components: {
        Target
    },
    watch: {
        itemBeingDragged: function(being_dragged) {
            if (being_dragged) {
                this.$refs['overlay'].classList.add('visible');
            } else {
                this.$refs['overlay'].classList.remove('visible');
            }
        },
    },
    computed: {
        ...mapState(['roomObjects', 'itemBeingDragged']),
        sortedRoomObjects: function() {
            return this.roomObjects.slice().sort((a, b) => {
                if (b.sort > a.sort) {
                    return 1;
                } else if (a.sort > b.sort) {
                    return -1;
                } else {
                    return 0;
                }
            })
        }
    },
    methods: {
        handleClick: function() {
            //this.$store.dispatch('setObjectTarget', '');
        },

        handleDragEnter: function() {
            this.$refs['overlay'].classList.add('highlight');
        },

        handleDragLeave: function() {
            this.$refs['overlay'].classList.remove('highlight');
        },

        handleDrop: function(e) {
            this.$refs['overlay'].classList.remove('highlight');
            let iuuid = e.dataTransfer.getData("item_uuid");
            this.$store.dispatch('sendSlashCommand', {
                command: `/drop ${iuuid}`
            });
        }
    }
}
</script>

<style lang="scss" scoped>
    .root {
        height: 100%;
        background-color: #131313;
    }

    .container {
        padding-top: 10px;
    }

    .drop-overlay {
        display: flex;
        opacity: 0;
        position: absolute;
        z-index: 10;
        width: 245px;
        height: 549px;
        left: -1000px;
        background-color: rgba(0, 0, 0, 0.72);
        align-items: center;
        justify-content: center;
        border: 2px dashed #353535;
        transition: opacity 0.1s ease-in-out;
        text-align: center;
        font-size: 18px;
        color: #777;
    }

    .drop-overlay.visible {
        left: 0px;
        opacity: 1;
    }

    .drop-overlay.highlight {
        border: 2px dashed #fff;
        color: #bbb;
    }
</style>