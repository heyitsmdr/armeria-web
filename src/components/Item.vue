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

<script setup>
    import { ref, computed, onMounted, inject } from 'vue';
    import { useStore } from 'vuex';
    import {INVENTORY_DRAG_START, INVENTORY_DRAG_STOP} from "@/plugins/SFX";

    const props = defineProps({
        uuid: String,
        name: String,
        slotNum: Number,
        equipSlot: String,
        pictureKey: String,
        color: String,
        equipped: Boolean
    });

    const item = ref(null); // Auto-mapped to HTML reference.

    const store = useStore();
    const isProduction = computed(() => store.state.isProduction);
    const itemTooltipUUID = computed(() => store.state.itemTooltipUUID);
    const itemTooltipVisible = computed(() => store.state.itemTooltipVisible);
    const itemTooltipMouseCoords = computed(() => store.state.itemTooltipMouseCoords);
    const hasPermission = computed(() => store.getters.hasPermission);

    const sfx = inject('sfx');

    onMounted(() => {
        item.value.classList.add('equipped');
    });

    /**
     * Handles when an item is being dragged over this item.
     * @param {DragEvent} e
     */
    function handleItemDragEnter(e) {
        e.target.classList.add('candrop');
    }

    /**
     * Handles when an item is no longer beign dragged over this item.
     * @param {DragEvent} e
     */
    function handleItemDragLeave(e) {
        e.target.classList.remove('candrop');
    }

    /**
     * Handles when the player starts dragging this item.
     * @param {DragEvent} e
     */
    function handleItemDragStart(e) {
        e.target.classList.add('dragging');
        e.dataTransfer.setData('item_uuid', props.uuid);
        e.dataTransfer.setData('item_slot', props.slotNum);
        hideTooltip();
        store.dispatch('setItemBeingDragged', true);
        store.dispatch('contextMenu/hide');
        sfx.play(INVENTORY_DRAG_START);
    }

    /**
     * Handles when the player stops dragging this item.
     * @param {DragEvent} e
     */
    function handleItemDragEnd(e) {
        e.target.classList.remove('dragging');
        store.dispatch('setItemBeingDragged', false);
        sfx.play(INVENTORY_DRAG_STOP);
    }

    /**
     * Handles when an item is dropped here. It is swapped.
     * @param {DragEvent} e
     */
    function handleItemDrop(e) {
        e.target.classList.remove('candrop');

        let slot = e.dataTransfer.getData("item_slot");
        if (slot) {
            store.dispatch('sendSlashCommand', {
                command: `/swap ${slot} ${props.slotNum}`,
                hidden: true,
            });
        }
    }

    /**
     * Handles when the mouse moves over the item (to show the tooltip).
     * @param {MouseEvent} e
     */
    function handleMouseMove(e) {
        if (!props.uuid) {
            return;
        }

        if (itemTooltipMouseCoords.value.x !== e.clientX || itemTooltipMouseCoords.value.y !== e.clientY) {
            store.dispatch('moveItemTooltip', { x: e.clientX, y: e.clientY });
        }

        if (itemTooltipUUID.value !== props.uuid) {
            store.dispatch('showItemTooltip', props.uuid);
        }
    }

    /**
     * Handles when the mouse leaves the item.
     */
    function handleMouseLeave() {
        hideTooltip();
    }

    /**
     * Handles the mouse up event on the item.
     * @param {MouseEvent} e
     */
    function handleMouseUp(e) {
        if (!props.uuid) {
            return;
        }

        if (e.shiftKey && hasPermission.value('CAN_BUILD')) {
            store.dispatch('sendSlashCommand', {
                command: `/item iedit ${props.uuid}`,
                hidden: true,
            });
        }
    }

    /**
     * Handles the right-click menu on the item.
     * @param {MouseEvent} e
     */
    function handleContextMenu(e) {
        if (!props.uuid) {
            return;
        }

        const items = [`Look %s|/look inv:${props.uuid}`];

        if (props.equipSlot.length > 0) {
            items.push(`Equip %s|/equip ${props.uuid}`);
        }

        items.push(
            `Wiki %s|wiki:/items/%s`,
            `Drop %s|/drop ${props.uuid}`,
            `Edit %s|/item iedit ${props.uuid}||CAN_BUILD`,
            `Edit-Parent %s|/item edit ${props.name}||CAN_BUILD`,
            `Destroy %s|/destroy ${props.uuid}||CAN_BUILD`,
        );

        store.dispatch(
            'contextMenu/show',
            {
                object: {
                    name: props.name,
                    color: `#${props.color}`,
                },
                at: {
                    x: e.pageX,
                    y: e.pageY,
                },
                items: items,
            }
        );
    }

    /**
     * Hides the tooltip, if it is visible.
     */
    function hideTooltip() {
        if (itemTooltipVisible.value) {
            store.dispatch('hideItemTooltip');
        }
    }

    /**
     * Returns a URL representing the picture for the item.
     * @returns {string}
     */
    function getBackgroundUrl() {
        if (!props.pictureKey) {
            return '';
        }

        if (!isProduction.value) {
            return `url(http://${window.location.hostname}:8081/oi/${props.pictureKey})`;
        }

        return `url(/oi/${props.pictureKey})`;
    }
</script>

<style>
    .tooltip .name {
        font-size: 20px;
        font-weight: 600;
    }
</style>
<style scoped lang="scss">
    @import "~@/styles/common";

    .item {
        width: 40px;
        height: 40px;
        background-color: $bg-color-light2;
        background-size: contain;
        margin: 2px;
        transition: all .1s ease-in-out;
        overflow: hidden;
        border: $defaultBorder;
        border-top-color: $bg-color;
        border-left-color: $bg-color;
        box-sizing: border-box;
    }

    .item:hover {
        cursor: pointer;
        border-color: $bg-color-light3;//transform: scale(1.1);
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
        background-color: $bg-color;
        padding: 5px;
    }

    .tooltip.visible {
        display: block;
    }
</style>