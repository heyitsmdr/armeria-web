<template>
    <div class="main-text-container" :style="{ height: containerHeight }">
        <ObjectEditor :style="{ height: containerHeight }"></ObjectEditor>
        <div class="scrollable-container" ref="mainTextContainer">
            <div class="lines">
                <div class="line" v-for="line in gameText" v-html="line"></div>
            </div>
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
                const height = this.windowHeight - 37 - 30 - 2;
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
    .main-text-container {
        display: flex;
    }

    .scrollable-container {
        padding: 5px;
        overflow-y: scroll;
        flex-grow: 1;
    }

    .line {
        color: #cacaca;
        user-select: all;
        white-space: pre;
    }
</style>