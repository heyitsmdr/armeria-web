<template>
    <div
        class="item"
        draggable="true"
        @dragstart="handleItemDragStart"
        @dragend="handleItemDragEnd"
        @dragenter="handleItemDragEnter"
        @dragleave="handleItemDragLeave"
    >

    </div>
</template>

<script>
export default {
    name: 'Item',
    methods: {
        handleItemDragEnter: function(e) {
            e.target.classList.add('candrop');
        },

        handleItemDragLeave: function(e) {
            e.target.classList.remove('candrop');
        },

        handleItemDragStart: function(e) {
            e.target.classList.add('dragging');
            //e.dataTransfer.setData('is_item', 'true');
            // And in Drop event: e.dataTransfer.getData("some_data");
        },

        handleItemDragEnd: function(e) {
            e.target.classList.remove('dragging');
        }
    }
}
</script>

<style scoped>
    .item {
        width: 40px;
        height: 40px;
        background-color: #000000;
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
</style>