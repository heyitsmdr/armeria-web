<template>
    <div class="container" ref="mainTextContainer" :style="{ height: containerHeight }">
        <ObjectEditor :style="{ height: containerHeight }"></ObjectEditor>
        <div class="lines">
            <div class="line" v-for="line in gameText" v-html="line"></div>
        </div>
    </div>
</template>

<script>
    import { mapState } from 'vuex'
    import ObjectEditor from "./ObjectEditor";

    export default {
        name: 'MainText',
        components: { ObjectEditor },
        data: function() {
            return {
                lineNumber: 0
            }
        },
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

<style scoped>
    .container {
        padding: 5px;
        overflow-y: scroll;
        box-shadow: inset 0px 0px 4px 0px #000;
        display: flex;
    }

    .lines {
        flex-grow: 1;
    }

    .line {
        color: #cacaca;
        user-select: all;
        white-space: pre;
    }
</style>