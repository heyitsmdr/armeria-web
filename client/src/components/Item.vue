<template>
    <div>
        <div
                class="item"
                draggable="true"
                :style="{ backgroundImage: getBackgroundUrl(), borderColor: color }"
                @dragstart="handleItemDragStart"
                @dragend="handleItemDragEnd"
                @dragenter="handleItemDragEnter"
                @dragleave="handleItemDragLeave"
                @dragover.prevent
                @drop="handleItemDrop"
                @mousemove="handleMouseMove"
                @mouseleave="handleMouseLeave"
        >
        </div>
        <div class="tooltip" ref="tooltip">
            {{ tooltipData }}
        </div>
    </div>
</template>

<script>
import { mapState } from 'vuex';
import {INVENTORY_DRAG_START, INVENTORY_DRAG_STOP} from "../plugins/SFX";

export default {
    name: 'Item',
    props: ['uuid', 'slotNum', 'pictureKey', 'tooltipData', 'color'],
    computed: mapState(['isProduction']),
    methods: {
        handleItemDragEnter: function(e) {
            e.target.classList.add('candrop');
        },

        handleItemDragLeave: function(e) {
            e.target.classList.remove('candrop');
        },

        handleItemDragStart: function(e) {
            e.target.classList.add('dragging');
            e.dataTransfer.setData('item_uuid', this.uuid);
            e.dataTransfer.setData('item_slot', this.slotNum);
            this.hideTooltip();
            this.$soundEvent(INVENTORY_DRAG_START);
        },

        handleItemDragEnd: function(e) {
            e.target.classList.remove('dragging');
            this.$soundEvent(INVENTORY_DRAG_STOP);
        },

        handleItemDrop: function(e) {
            e.target.classList.remove('candrop');

            let slot = e.dataTransfer.getData("item_slot");
            if (slot) {
                this.$store.dispatch('sendSlashCommand', {
                    command: `/swap ${slot} ${this.slotNum}`
                });
            }
        },

        handleMouseMove: function(e) {
            if (!this.uuid) {
                return;
            }

            const xOffset = 50;
            const yOffset = 50;

            let tt = this.$refs["tooltip"];
            if (!tt.classList.contains('visible')) {
                tt.classList.add('visible');
            }

            let ttTop = e.clientY - tt.clientHeight - yOffset;
            let ttLeft = e.clientX - xOffset;

            if ((ttLeft + tt.clientWidth + 20) > window.innerWidth) {
                ttLeft = window.innerWidth - tt.clientWidth - 20;
            }

            tt.style.top = ttTop + 'px';
            tt.style.left = ttLeft + 'px';
            tt.style.borderColor = this.color;
        },

        handleMouseLeave: function() {
            this.hideTooltip();
        },

        hideTooltip: function() {
            let tt = this.$refs["tooltip"];
            tt.classList.remove('visible');
        },

        getBackgroundUrl() {
            if (!this.pictureKey) {
                return '';
            }

            if (!this.isProduction) {
                return `url(http://localhost:8081/oi/${this.pictureKey})`;
            }

            return `url(/oi/${this.pictureKey})`;
        },
    }
}
</script>

<style scoped>
    .item {
        width: 40px;
        height: 40px;
        background-color: #000000;
        background-size: contain;
        border: 1px solid #333;
        margin-right: 4px;
        margin-bottom: 4px;
        transition: all .1s ease-in-out;
        overflow: hidden;
    }

    .item:hover {
        cursor: pointer;
        transform: scale(1.1);
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

    .tooltip {
        display: none;
        position: absolute;
        width: 200px;
        z-index: 999;
        top: 50px;
        background-color: #111;
        border: 2px solid #ccc;
        padding: 5px;
        box-shadow: 0px 0px 10px #000;
    }

    .tooltip.visible {
        display: block;
    }
</style>