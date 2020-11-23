<template>
    <div
        ref="menu"
        class="menu"
        :style="{ top: `${this.contextMenuPosition.y + 5}px` }"
    >
        <div
            class="item"
            v-for="item in contextMenuItems"
            :key="item"
            :class="itemClasses(item)"
            @click="handleItemClick(item)"
            @mouseenter="handleItemMouseEnter"
        >
            <span v-html="itemNameHTML(item)"></span>
        </div>
    </div>
</template>

<script>
    import {mapState} from 'vuex';
    import {INVENTORY_DRAG_START, INVENTORY_DRAG_STOP} from "@/plugins/SFX";
    export default {
        name: 'ContextMenu',
        mounted: function() {
            window.addEventListener('click', this.handleWindowClick);
        },
        watch: {
            contextMenuItems: function() {
                this.$soundEvent(INVENTORY_DRAG_START);

                this.$nextTick(() => {
                    const menu = this.$refs["menu"];
                    const width = menu.clientWidth;
                    const windowWidth = window.innerWidth;

                    if ((this.contextMenuPosition.x + width + 20) > windowWidth) {
                        menu.style.left = `${this.contextMenuPosition.x + 5 - width}px`;
                        return
                    }

                    menu.style.left = `${this.contextMenuPosition.x + 5}px`;
                });
            },

            contextMenuVisible: function(visible) {
                if (!visible) {
                    const menu = this.$refs["menu"];
                    menu.style.left = null;
                }
            }
        },
        computed: {
            ...mapState([
                'contextMenuVisible',
                'contextMenuItems',
                'contextMenuObjectName',
                'contextMenuObjectColor',
                'contextMenuPosition',
            ])
        },
        methods: {
            itemNameHTML: function(item) {
                const itemSections = item.split('|');
                let displayText = itemSections[0];
                let displayStyle = itemSections.length >= 3 ? itemSections[2] : '';

                displayText = displayText.replace('%s', `<span style="color:${this.contextMenuObjectColor};font-weight:600">[${this.contextMenuObjectName}]</span>`);

                if (displayStyle.length > 0) {
                    return `<span style="${displayStyle}">${displayText}</span>`;
                }

                return displayText;
            },

            itemClasses: function(item) {
                const itemSections = item.split('|');
                return itemSections.length >= 4 ? itemSections[3] : '';
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

                this.$soundEvent(INVENTORY_DRAG_STOP);
            },

            handleItemMouseEnter: function() {
                this.$soundEvent(INVENTORY_DRAG_START);
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
    @import "@/styles/common";
    .menu {
        position: absolute;
        z-index: 900;
        top: 0px;
        left: -500px;
        background-color: $defaultBackgroundColor;
        font-size: 13px;
    }

    .item {
        padding: 5px 12px;

        &.admin {
            background-color: #61030369;

            &:hover {
                background-color: #a9070769;
            }
        }

        &:not(:last-child) {
            border-bottom: 1px solid #4e4e4e;
        }

        &:hover {
            background-color: #403e3ee8;
            cursor: pointer;
        }
    }

</style>