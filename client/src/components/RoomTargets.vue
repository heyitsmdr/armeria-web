<template>
    <div class="root">
        <div class="container">
            <Target v-for="obj in roomObjects" :name="obj.name" />
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
    computed: mapState(['roomObjects'])
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