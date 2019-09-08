<template>
    <div class="root">
        <div class="container" @click="handleClick">
            <Target
                v-for="obj in sortedRoomObjects"
                :key="obj.uuid"
                :uuid="obj.uuid"
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