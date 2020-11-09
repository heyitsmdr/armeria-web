<template>
    <div class="targets-container">
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

<style scoped>
.targets-container {
    display: flex;
    align-items: center;
    height: 100%;
    margin: 0 5px;
    font-size: 12px;
}

.ping {

}
</style>