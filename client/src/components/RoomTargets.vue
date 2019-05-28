<template>
    <div class="root">
        <div class="container" @click="handleClick">
            <Target v-for="obj in roomObjects" :name="obj.name" :pictureKey="obj.picture" />
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
    computed: mapState(['roomObjects']),
    methods: {
        handleClick: function() {
            //this.$store.dispatch('setObjectTarget', '');
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
</style>