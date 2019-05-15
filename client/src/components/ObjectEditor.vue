<template>
    <div class="object-editor" :class="{ 'visible': objectEditorOpen }">
        This is the object editor.
    </div>
</template>

<script>
    import { mapState } from 'vuex';

    export default {
        name: 'ObjectEditor',
        computed: mapState(['objectEditorOpen'])
    }
</script>

<style scoped>
    .object-editor {
        background-color: #0b0b0b;
        flex-basis: 0px;
        overflow: hidden;
        transition: all .2s ease-in-out;
    }

    .object-editor.visible {
        flex-basis: 300px;
        margin-right: 5px;
    }
</style>