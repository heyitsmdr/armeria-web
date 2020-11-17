<template>
    <div
        ref="menu"
        class="menu"
        :style="{ top: `${this.contextMenuPosition.y + 5}px`, left: this.calcLeft }"
    >
        <div
            class="item"
            v-for="item in contextMenuItems"
            :key="item"
            @click="handleItemClick(item)"
        >
            <span v-html="itemNameHTML(item)"></span>
        </div>
    </div>
</template>

<script>
    import {mapState} from 'vuex';
    export default {
        name: 'ContextMenu',
        mounted: function() {
            window.addEventListener('click', this.handleWindowClick);
        },
        computed: {
            calcLeft: function() {
                if (!this.contextMenuVisible) {
                    return '-500px';
                } else if (this.$refs["menu"]) {
                    const width = this.$refs["menu"].clientWidth;
                    const windowWidth = window.innerWidth;

                    // If the context menu would appear off-screen, move it to the left of the cursor position.
                    if ((this.contextMenuPosition.x + width + 20) > windowWidth) {
                        return `${this.contextMenuPosition.x + 5 - width}px`;
                    }
                }

                return `${this.contextMenuPosition.x + 5}px`;
            },
            ...mapState([
                'contextMenuVisible',
                'contextMenuItems',
                'contextMenuObjectName',
                'contextMenuObjectColor',
                'contextMenuPosition',
            ])
        },
        watch: {
        },
        methods: {
            itemNameHTML: function(item) {
                let displayText = item.split('|')[0];
                displayText = displayText.replace('%s', `<span style="color:${this.contextMenuObjectColor};font-weight:600">[${this.contextMenuObjectName}]</span>`);
                return displayText;
            },

            handleItemClick: function(item) {
                const cmd = item.split('|')[1];

                if (cmd.substr(0, 5) === 'wiki:') {
                    let wikiUrl = cmd.substr(5);
                    wikiUrl = wikiUrl.replace('%s', this.contextMenuObjectName.toLowerCase().replaceAll(' ', '-'));
                    window.open(`https://wiki.armeria.io${wikiUrl}`);
                    return;
                }

                this.$store.dispatch('sendSlashCommand', {
                    command: cmd
                });

                this.$store.dispatch('hideContextMenu');
            },

            handleWindowClick: function() {
                if (this.contextMenuVisible) {
                    this.$store.dispatch('hideContextMenu');
                }
            },
        }
    }
</script>

<style scoped lang="scss">
    .menu {
        position: absolute;
        z-index: 900;
        top: 200px;
        left: 800px;
        background-color: #2d2c2c;
        border: 1px solid #4e4e4e;
    }

    .item {
        padding: 5px 12px;

        &:not(:last-child) {
            border-bottom: 1px solid #4e4e4e;
        }

        &:hover {
            background-color: #4e4e4e;
            cursor: pointer;
        }
    }

</style>