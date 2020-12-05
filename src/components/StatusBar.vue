<template>
    <div class="root">
        <div class="inner-container">Ping: {{pingTime}}ms</div>
    </div>
</template>

<script>
    import {mapState} from "vuex";

    export default {
        name: 'StatusBar',
        computed: mapState(['pingTime']),
    }
</script>

<style scoped lang="scss">
    @import "@/styles/common";

    .root {
        display: flex;
        height: 100%;
        align-items: center;
        font-size: 12px;
        box-sizing: border-box;
    }

    .ping {

    }
</style>