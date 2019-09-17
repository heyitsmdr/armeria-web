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

<style scoped>
    .root {
        height: 100%;
        background-color: #131313;
    }

    .container {
        padding-top: 10px;
    }
</style>