<template>
    <div class="root" :style="{ height: containerHeight }">
        <ObjectEditor :style="{ height: containerHeight }"></ObjectEditor>
        <div class="scrollable-container" ref="mainTextContainer">
            <div class="lines">
                <div class="line" v-for="line in gameText" v-html="line.html" :key="line.id"></div>
            </div>
        </div>
        <div
                class="item-drag-overlay"
                ref="item-overlay"
                @dragenter="handleItemOverlayDragEnter"
                @dragleave="handleItemOverlayDragLeave"
                @drop="handleItemOverlayDrop"
                @dragover.prevent
                v-if="itemBeingDragged"
        >
            Release the item here to drop it into the room.
        </div>
    </div>
</template>

<script>
import {mapGetters, mapState} from 'vuex'
    import ObjectEditor from "./ObjectEditor";

    export default {
        name: 'MainText',
        components: {ObjectEditor},
        data: function () {
            return {
                lineNumber: 0,
                lastItemTooltipUUID: '',
            }
        },
        props: {
            windowHeight: Number,
        },
        computed: {
            ...mapState(['gameText', 'itemBeingDragged', 'settings']),
            ...mapGetters(['hasPermission']),
            containerHeight() {
                const height = this.windowHeight - 37 - 30 - 2 - 35;
                return `${height}px`;
            }
        },
        watch: {
            gameText: function(lines) {
                let maxLines = this.settings['lines'];
                if (!maxLines) {
                    return;
                }

                maxLines = parseInt(maxLines);

                if (maxLines > lines.length) {
                    // Delete the oldest line here.
                }
            }
        },
        updated: function () {
            this.$nextTick(function () {
                const div = this.$refs['mainTextContainer'];
                div.scrollTop = 9999999;
            });
        },
        mounted: function () {
            document.addEventListener('click', e => {
                if (e.target.className === 'convo-select') {
                    const optionDisabled = e.target.getAttribute('data-disabled');
                    const groupId = e.target.getAttribute('data-group');
                    const convoOptionId = e.target.getAttribute('data-convo-option-id');
                    const mobUUID = e.target.getAttribute('data-mob-uuid');

                    if (optionDisabled && optionDisabled === 'true') {
                        return;
                    }

                    const groupSpans = document.querySelectorAll(`.convo-select[data-group="${groupId}"]`);
                    groupSpans.forEach(span => {
                        span.style.color = '#444';
                        span.style.borderLeftColor = '#444';
                        span.style.borderBottomColor = '#292929';
                        span.setAttribute('data-disabled', 'true');
                    });

                    this.$store.dispatch('sendSlashCommand', {
                        command: `/select "${mobUUID}" "${convoOptionId}"`,
                        hidden: true,
                    });
                } else if (e.target.className === 'inline-command') {
                    const commandEncoded = e.target.getAttribute('data-command');
                    const buff = new Buffer(commandEncoded, 'base64');
                    this.$store.dispatch('sendSlashCommand', {
                        command: buff.toString('ascii'),
                        hidden: true,
                    });
                }
            });

            document.addEventListener('mousemove', e => {
                if (e.target.className === 'hover-item-tooltip') {
                    const uuid = e.target.getAttribute('data-uuid');
                    if (this.lastItemTooltipUUID !== uuid) {
                        this.lastItemTooltipUUID = uuid;
                        this.$store.dispatch('showItemTooltip', uuid);
                    }
                    this.$store.dispatch('moveItemTooltip', { x: e.clientX, y: e.clientY });
                } else if (this.lastItemTooltipUUID.length > 0) {
                    this.lastItemTooltipUUID = '';
                    this.$store.dispatch('hideItemTooltip');
                }
            }, false);

            document.addEventListener('contextmenu', e => {
                let menuSpan;
                if (e.target.className === 'hover-item-tooltip') {
                  if (e.path[1].className === 'dynamic-context-menu') {
                      menuSpan = e.path[1];
                  }
                } else if (e.target.className === 'dynamic-context-menu') {
                    menuSpan = e.target;
                }

                if (!menuSpan) {
                    return;
                }

                e.preventDefault();
                e.stopPropagation();

                const contentEncoded = menuSpan.getAttribute('data-content');
                const buff = new Buffer(contentEncoded, 'base64');
                let menuItems = buff.toString('ascii').replaceAll('@', '%s').split(';');
                menuItems = menuItems.filter(c => {
                    const sections = c.split('|');
                    if (sections.length >= 4 && sections[3] === 'admin') {
                        return this.hasPermission('CAN_BUILD');
                    }

                    return true;
                });

                this.$store.dispatch(
                    'showContextMenu',
                    {
                        object: {
                            name: menuSpan.getAttribute('data-name'),
                            color: `#${menuSpan.getAttribute('data-color')}`,
                            subjectBrackets: (menuSpan.getAttribute('data-type') === 'item'),
                        },
                        at: {
                            x: e.pageX,
                            y: e.pageY,
                        },
                        items: menuItems,
                    }
                );
            }, false);
        },
        methods: {
            handleItemOverlayDragEnter: function () {
                this.$refs['item-overlay'].classList.add('item-over');
            },

            handleItemOverlayDragLeave: function () {
                this.$refs['item-overlay'].classList.remove('item-over');
            },

            handleItemOverlayDrop: function (e) {
                this.$refs['item-overlay'].classList.remove('item-over');
                let iuuid = e.dataTransfer.getData("item_uuid");
                this.$store.dispatch('sendSlashCommand', {
                    command: `/drop "${iuuid}"`,
                    hidden: true,
                });
            }
        },
    }
</script>

<style>
    .line .monospace {
        font-family: 'Inconsolata', monospace;
        font-size: 16px;
        white-space: pre;
    }

    .line .inline-link {
        color: #666;
    }

    .line .inline-link:hover {
        color: #aaa;
    }

    .line .inline-command {
        color: #2196F3;
        text-decoration: none;
        font-weight: 600;
    }

    .line .inline-command:hover {
        color: #00BCD4;
    }

    .line .inline-button {
        background-color: #383737;
        padding: 0px 5px;
        border: 1px solid #585555;
    }

    .line .inline-button:hover {
        cursor: pointer;
        border: 1px solid #848282;
    }

    .line .inline-loopback {
        color: #5d5d5d;
        padding-top: 15px;
    }

    .line .hover-item-tooltip:hover {
        cursor: pointer;
        border-bottom: 1px dotted #666;
    }

    .line .dynamic-context-menu:hover {
        cursor: pointer;
        border-bottom: 1px dotted #666;
    }

    .line .convo-select {
        border-left: 2px solid #86949e;
        padding-left: 5px;
        margin-left: 10px;
        border-bottom: 1px solid #2b3740;
        padding-bottom: 2px;
        transition: all .1s ease-in-out;
    }

    .line .convo-select:hover {
        cursor: pointer;
        border-left-width: 4px;
        color: #fff;

    }

    .line table {
        border: 2px solid #1d2e39;
    }

    .line table tr th {
        text-align: left;
        background: linear-gradient(to bottom, #1d2e39 0%, #132029 100%);
        border-bottom: 2px solid #1d2e39;
        padding: 3px;
    }

    .line table tr {
        background-color: #132029;
    }

    .line table tr:hover {
        background-color: #0d181f;
    }
    .line table tr td {
        padding: 2px 3px;
    }
</style>
<style scoped lang="scss">
    @import "~@/styles/common";
    
    .root {
        box-sizing: border-box;
        display: flex;

        /*border: $defaultBorder;*/
        @include defaultBorderImage;
    }

    .root .item-drag-overlay {
        position: absolute;
        z-index: 100;
        top: 0;
        bottom: 0;
        background-color: #000000b8;
        width: 100%;
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 22px;
        color: #666;
        border: 2px dashed #666;
    }

    .root .item-drag-overlay.item-over {
        background-color: #1d1c1cb8;
        border: 2px dashed #aaa;
        color: #aaa;
    }

    .scrollable-container {
        padding: 5px;
        overflow-y: scroll;
        flex-grow: 1;
    }

    .line {
        color: #cacaca;
        user-select: text;
        font-size: 13px;
    }
</style>