<template>
    <div class="container">
        <div class="line" v-for="line in gameText">
            {{ line }}
        </div>
    </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
    name: 'MainText',
    computed: mapState(['gameText'])
}
</script>

<style lang="scss" scoped>
    .container {
        padding: 5px;
    }
</style>