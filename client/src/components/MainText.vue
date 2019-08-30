<template>
    <div class="main-text-container" :style="{ height: containerHeight }">
        <ObjectEditor :style="{ height: containerHeight }"></ObjectEditor>
        <div class="scrollable-container" ref="mainTextContainer">
            <div class="lines">
                <div class="line" v-for="line in gameText" v-html="line.html" :key="line.id"></div>
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

<style>
    .line .monospace {
        font-family: 'Inconsolata', monospace;
        font-size: 16px;
        white-space: pre;
    }

    .line table tr th {
        text-align: left;
        background: linear-gradient(to bottom, #111111 0%,#232323 100%);
        padding: 3px;
    }

    .line table tr td {
        padding: 3px;
    }
</style>
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
        user-select: text;
    }
</style>