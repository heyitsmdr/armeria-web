<template>
    <div
        ref="menu"
        class="menu"
        :style="{
            top: `${this.contextMenuPosition.y + 5}px`,
            opacity: (this.contextMenuVisible) ? 1 : 0,
        }"
        :class="{
            visible: this.contextMenuVisible,
        }"
    >
        <div
            class="item"
            v-for="item in filteredMenuItems"
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
import {mapGetters, mapState} from 'vuex';
    import {INVENTORY_DRAG_START, INVENTORY_DRAG_STOP} from "@/plugins/SFX";
    export default {
        name: 'ContextMenu',
        mounted: function() {
            window.addEventListener('click', this.handleWindowClick);
        },
        data: function() {
            return {
                hideTimeout: null,
            }
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
                    // The timeout (200s) should match the transition duration of the .menu class.
                    this.hideTimeout = setTimeout(() => {
                        const menu = this.$refs["menu"];
                        menu.style.left = null;
                        this.hideTimeout = null;
                    }, 200);
                } else if (this.hideTimeout !== null) {
                    clearTimeout(this.hideTimeout);
                    this.hideTimeout = null;
                }
            }
        },
        computed: {
            ...mapState([
                'contextMenuVisible',
                'contextMenuItems',
                'contextMenuObjectName',
                'contextMenuObjectColor',
                'contextMenuObjectBrackets',
                'contextMenuPosition',
            ]),
            ...mapGetters(['hasPermission']),
            filteredMenuItems: function() {
                return this.contextMenuItems.filter(item => {
                    const sections = item.split('|');
                    if (sections.length >= 4) {
                        return this.hasPermission(sections[3]);
                    }

                    return true;
                });
            },
        },
        methods: {
            itemNameHTML: function(item) {
                const itemSections = item.split('|');
                let displayText = itemSections[0];
                let displayStyle = itemSections.length >= 3 ? itemSections[2] : '';

                if (this.contextMenuObjectBrackets) {
                    displayText = displayText.replace('%s', `<span style="color:${this.contextMenuObjectColor};font-weight:600">[${this.contextMenuObjectName}]</span>`);
                } else {
                    displayText = displayText.replace('%s', `<span style="color:${this.contextMenuObjectColor};font-weight:600">${this.contextMenuObjectName}</span>`);
                }

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
                    command: cmd,
                    hidden: true,
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
    @import "~@/styles/common";
    $borderColor: #4e4e4e;

    .menu {
        position: absolute;
        z-index: 900;
        top: 0px;
        left: -500px;
        background-color: $bg-color;
        font-size: 12px;
        border: 2px solid $borderColor;
        box-shadow: 0px 0px 10px #000;
        transition: opacity .2s ease-in-out;

        &.visible {
            max-height: 500px;
        }
    }

    .item {
        padding: 5px 12px;

        &.CAN_BUILD, &.CAN_CHAREDIT {
            background-color: #61030369;

            &:hover {
                background-color: #a9070769;
            }
        }

        &:not(:last-child) {
            border-bottom: 1px solid $borderColor;
        }

        &:hover {
            background-color: #403e3ee8;
            cursor: pointer;
        }
    }

</style>