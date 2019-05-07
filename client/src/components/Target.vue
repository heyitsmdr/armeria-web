<template>
    <div>
        <div class="container" @mousedown="onMouseOver" @mouseup="onMouseLeave">
            {{ name }}
        </div>
    </div>
</template>

<script>
export default {
    name: 'Target',
    props: ['name'],
    methods: {
        onMouseOver: function(event) {
            event.target.classList.add('mouse-down');
        },

        onMouseLeave: function(event) {
            event.target.classList.remove('mouse-down');
        }
    }
}
</script>

<style lang="scss" scoped>
.container {
    background-color: #0c0c0c;
    margin: 0 10px 10px 10px;
    padding: 15px;
    transition: all .1s ease-in-out;

    &.mouse-down {
        transform: scale(1.01) !important;
    }

    &:hover {
        cursor: pointer;
        transform: scale(1.05);
    }
}
</style>