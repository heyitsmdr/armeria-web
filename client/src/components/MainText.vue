<template>
    <div class="container" ref="mainTextContainer" v-bind:style="{ height: containerHeight }">
        <div class="line" v-for="line in gameText" v-html="line"></div>
    </div>
</template>

<script>
import { mapState } from 'vuex'

export default {
    name: 'MainText',
    props: {
        windowHeight: Number,
    },
    computed: {
        ...mapState(['gameText']),
        containerHeight() {
            const height = this.windowHeight - 300 - 45;
            return `${height}px`;
        }
    },
    updated: function() {
        this.$nextTick(function () {
            const div = this.$refs['mainTextContainer'];
            div.scrollTop = 9999999;
        });
    },
    methods: {

    },
    mounted: function() {

    },
    destroyed: function() {

    }
}
</script>

<style lang="scss" scoped>
    .container {
        padding: 5px;
        overflow-y: scroll;
    }

    .line {
        color: #cacaca;
    }
</style>