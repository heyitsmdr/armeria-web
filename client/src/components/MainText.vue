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
}
</script>

<style lang="scss" scoped>
    .container {
        padding: 5px;
        overflow-y: scroll;
        box-shadow: inset 0px 0px 4px 0px #000;
    }

    .line {
        color: #cacaca;
        user-select: all;
    }
</style>