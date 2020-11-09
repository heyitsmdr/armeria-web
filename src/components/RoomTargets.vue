<template>
    <div class="root">
        <div class="banner">
            Room Targets
        </div>
        <div class="targets-container" @click="handleClick">
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

<style scoped>
    .root {
        height: 100%;
        background-color: #131313;
    }

    .banner {
        background-color: #1b1b1b;
        font-weight: 800;
        text-transform: uppercase;
        padding: 4px 5px;
        border-bottom: 1px solid #333;
        text-align: center;
        color: #b7b7b7;
        position: fixed;
        width: 236px;
        z-index: 10;
    }

    .targets-container {
        padding-top: 37px;
    }
</style>