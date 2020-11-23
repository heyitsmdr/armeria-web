<template>
    <div>
        <div
                class="item"
                ref="item"
                draggable="true"
                :style="{ backgroundImage: getBackgroundUrl(), borderColor: color ? `#${color}` : '' }"
                @dragstart="handleItemDragStart"
                @dragend="handleItemDragEnd"
                @dragenter="handleItemDragEnter"
                @dragleave="handleItemDragLeave"
                @dragover.prevent
                @drop="handleItemDrop"
                @mousemove="handleMouseMove"
                @mouseleave="handleMouseLeave"
                @mouseup="handleMouseUp"
                @contextmenu.stop.prevent="handleContextMenu"
        >
            <div v-if="equipped" class="equipped">equip</div>
        </div>
    </div>
</template>

<script>
    import {mapState, mapGetters} from 'vuex';
    import {INVENTORY_DRAG_START, INVENTORY_DRAG_STOP} from "../plugins/SFX";

    export default {
        name: 'Item',
        props: ['uuid', 'name', 'slotNum', 'pictureKey', 'color', 'equipped'],
        computed: {
            ...mapState(['isProduction', 'itemTooltipUUID', 'itemTooltipVisible', 'itemTooltipMouseCoords']),
            ...mapGetters(['hasPermission']),
        },
        mounted: function () {
            this.$refs['item'].classList.add('equipped');
        },
        methods: {
            handleItemDragEnter: function (e) {
                e.target.classList.add('candrop');
            },

            handleItemDragLeave: function (e) {
                e.target.classList.remove('candrop');
            },

            handleItemDragStart: function (e) {
                e.target.classList.add('dragging');
                e.dataTransfer.setData('item_uuid', this.uuid);
                e.dataTransfer.setData('item_slot', this.slotNum);
                this.hideTooltip();
                this.$store.dispatch('setItemBeingDragged', true);
                this.$store.dispatch('hideContextMenu');
                this.$soundEvent(INVENTORY_DRAG_START);
            },

            handleItemDragEnd: function (e) {
                e.target.classList.remove('dragging');
                this.$store.dispatch('setItemBeingDragged', false);
                this.$soundEvent(INVENTORY_DRAG_STOP);
            },

            handleItemDrop: function (e) {
                e.target.classList.remove('candrop');

                let slot = e.dataTransfer.getData("item_slot");
                if (slot) {
                    this.$store.dispatch('sendSlashCommand', {
                        command: `/swap ${slot} ${this.slotNum}`
                    });
                }
            },

            handleMouseMove: function (e) {
                if (!this.uuid) {
                    return;
                }

                if (this.itemTooltipMouseCoords.x !== e.clientX || this.itemTooltipMouseCoords.y !== e.clientY) {
                    this.$store.dispatch('moveItemTooltip', { x: e.clientX, y: e.clientY });
                }

                if (this.itemTooltipUUID !== this.uuid) {
                    this.$store.dispatch('showItemTooltip', this.uuid);
                }
            },

            handleMouseLeave: function () {
                this.hideTooltip();
            },

            handleMouseUp: function (e) {
                if (!this.uuid) {
                    return;
                }

                if (e.shiftKey && this.hasPermission('CAN_BUILD')) {
                    this.$socket.sendObj({
                        type: 'command',
                        payload: `/item iedit ${this.uuid}`
                    });
                }
            },

            handleContextMenu: function (e) {
                if (!this.uuid) {
                    return;
                }

                const items = [
                    `Look %s|/look inv:${this.uuid}`,
                    `Wiki %s|wiki:/items/%s`,
                    `Drop %s|/drop ${this.uuid}`,
                ];

                if (this.hasPermission('CAN_BUILD')) {
                    items.push(`Destroy %s|/destroy ${this.uuid}||admin`);
                    items.push(`Edit %s|/item iedit ${this.uuid}||admin`);
                    items.push(`Edit-Parent %s|/item edit ${this.name}||admin`);
                }
                this.$store.dispatch(
                    'showContextMenu',
                    {
                        object: {
                            name: this.name,
                            color: `#${this.color}`,
                        },
                        at: {
                            x: e.pageX,
                            y: e.pageY,
                        },
                        items: items,
                    }
                );


            },

            hideTooltip: function () {
                if (this.itemTooltipVisible) {
                    this.$store.dispatch('hideItemTooltip');
                }
            },

            getBackgroundUrl() {
                if (!this.pictureKey) {
                    return '';
                }

                if (!this.isProduction) {
                    return `url(http://${window.location.hostname}:8081/oi/${this.pictureKey})`;
                }

                return `url(/oi/${this.pictureKey})`;
            },
        }
    }
</script>

<style>
    .tooltip .name {
        font-size: 20px;
        font-weight: 600;
    }
</style>
<style scoped lang="scss">
    @import "@/styles/common";

    .item {
        width: 40px;
        height: 40px;
        background-color: $defaultBackgroundColor;
        background-size: contain;
        margin-right: 4px;
        margin-bottom: 4px;
        transition: all .1s ease-in-out;
        overflow: hidden;
        border: $defaultBorder;
        box-sizing: border-box;
    }

    .item:hover {
        cursor: pointer;
        transform: scale(1.1);
    }

    .item.equipped {

    }

    .item.candrop {
        transform: scale(1.4);
    }

    .item.dragging {
        opacity: 0.5;
    }

    .item .picture {
        background-size: contain;
        height: 100%;
        width: 100%;
    }

    .item .equipped {
        background-color: rgba(50, 50, 50, 0.8);
        color: #fff;
        font-size: 10px;
        text-align: center;
        margin-top: 27px;
        text-transform: uppercase;
    }

    .tooltip {
        display: none;
        position: absolute;
        max-width: 400px;
        min-width: 150px;
        z-index: 999;
        top: 50px;
        background-color: $defaultBackgroundColor;
        padding: 5px;
    }

    .tooltip.visible {
        display: block;
    }
</style>