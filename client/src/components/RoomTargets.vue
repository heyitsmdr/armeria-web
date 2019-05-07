<template>
    <div class="root">
        <div class="container">
            <Target name="Ethryx" />
            <Target name="OtherDude" />
        </div>
    </div>
</template>

<script>
import Target from '@/components/Target';

export default {
    name: 'RoomTargets',
    components: {
        Target
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