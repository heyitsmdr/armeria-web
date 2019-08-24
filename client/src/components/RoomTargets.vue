<template>
    <div class="root" @dragenter="handleDragEnter" @dragleave="handleDragLeave" @dragover.prevent @drop="handleDrop">
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
    computed: {
        ...mapState(['roomObjects']),
        sortedRoomObjects: function() {
            return this.roomObjects.sort((a, b) => {
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

        handleDragEnter: function(e) {
            e.target.classList.add('candrop');
        },

        handleDragLeave: function(e) {
            e.target.classList.remove('candrop');
        },

        handleDrop: function(e) {
            e.target.classList.remove('candrop');
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

    .root.candrop {
        background-color: #313131;
    }

    .container {
        padding-top: 10px;
    }

</style>