<template>
    <div
        ref="menu"
        class="menu"
        :style="{
            top: `${contextMenuPosition.y + 5}px`,
            opacity: (contextMenuVisible) ? 1 : 0,
        }"
        :class="{
            visible: contextMenuVisible,
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

<script setup>
    import { ref, computed, watch, inject, nextTick, onMounted } from 'vue';
    import { useStore } from "vuex";

    import { INVENTORY_DRAG_START, INVENTORY_DRAG_STOP } from "@/plugins/SFX";

    const sfx = inject('sfx');

    const hideTimeout = ref(null);
    const menu = ref(null); // Auto-mapped to HTML reference.

    const store = useStore();
    const contextMenuVisible = computed(() => store.state.contextMenu.visible);
    const contextMenuItems = computed(() => store.state.contextMenu.items);
    const contextMenuObjectName = computed(() => store.state.contextMenu.objectName);
    const contextMenuObjectColor = computed(() => store.state.contextMenu.objectColor);
    const contextMenuObjectBrackets = computed(() => store.state.contextMenu.objectBrackets);
    const contextMenuPosition = computed(() => store.state.contextMenu.position);
    const hasPermission = computed(() => store.getters.hasPermission);

    const filteredMenuItems = computed(() => {
        return contextMenuItems.value.filter(item => {
            const sections = item.split('|');
            if (sections.length >= 4) {
                return hasPermission.value(sections[3]);
            }

            return true;
        });
    });

    watch(contextMenuItems, async () => {
        sfx.play(INVENTORY_DRAG_START);

        await nextTick();

        const width = menu.value.clientWidth;
        const windowWidth = window.innerWidth;

        if ((contextMenuPosition.value.x + width + 20) > windowWidth) {
            menu.value.style.left = `${contextMenuPosition.value.x + 5 - width}px`;
            return
        }

        menu.value.style.left = `${contextMenuPosition.value.x + 5}px`;
    });

    watch(contextMenuVisible, (visible) => {
        if (!visible) {
            // The timeout (200s) should match the transition duration of the .menu class.
            hideTimeout.value = setTimeout(() => {
                menu.value.style.left = null;
                hideTimeout.value = null;
            }, 200);
        } else if (hideTimeout.value !== null) {
            clearTimeout(hideTimeout.value);
            hideTimeout.value = null;
        }
    });

    onMounted(() => {
        window.addEventListener('click', handleWindowClick);
    });

    /**
     * Returns the HTML-version of the item name.
     * @param {String} item
     * @returns {String}
     */
    function itemNameHTML(item) {
        const itemSections = item.split('|');
        let displayText = itemSections[0];
        let displayStyle = itemSections.length >= 3 ? itemSections[2] : '';

        if (contextMenuObjectBrackets.value) {
            displayText = displayText.replace('%s', `<span style="color:${contextMenuObjectColor.value};font-weight:600">[${contextMenuObjectName.value}]</span>`);
        } else {
            displayText = displayText.replace('%s', `<span style="color:${contextMenuObjectColor.value};font-weight:600">${contextMenuObjectName.value}</span>`);
        }

        if (displayStyle.length > 0) {
            return `<span style="${displayStyle}">${displayText}</span>`;
        }

        return displayText;
    }

    /**
     * Returns any special classes for the menu entry.
     * @param {String} item
     * @returns {String}
     */
    function itemClasses(item) {
        const itemSections = item.split('|');
        return itemSections.length >= 4 ? itemSections[3] : '';
    }

    /**
     * Handles clicking the menu item.
     * @param {String} item
     */
    function handleItemClick(item) {
        const cmd = item.split('|')[1];

        if (cmd.substring(0, 5) === 'wiki:') {
            let wikiUrl = cmd.substring(5);
            wikiUrl = wikiUrl.replace('%s', contextMenuObjectName.value.toLowerCase().replaceAll(' ', '-'));
            window.open(`https://wiki.armeria.io${wikiUrl}`);
            return;
        }

        store.dispatch('sendSlashCommand', {
            command: cmd,
            hidden: true,
        });

        store.dispatch('contextMenu/hide');

        sfx.play(INVENTORY_DRAG_STOP);
    }

    /**
     * Handles the mouse hovering over the menu item.
     */
    function handleItemMouseEnter() {
        sfx.play(INVENTORY_DRAG_START);
    }

    /**
     * Handles clicking on the window, to hide the context menu.
     */
    function handleWindowClick() {
        if (contextMenuVisible.value) {
            store.dispatch('contextMenu/hide');
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