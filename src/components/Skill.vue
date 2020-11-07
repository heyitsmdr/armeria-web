<template>
    <div class="skill-container">
        <div class="name-container">
            <div class="name">{{ name }}</div>
            <div class="level">{{ level }}</div>
        </div>
        <div class="bar-container">
            <div class="bar" ref="bar"></div>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'Skill',
        props: ['name', 'level', 'percent'],
        mounted: function() {
            this.$refs['bar'].style.width = `${this.percent}%`;
        },
    }
</script>

<style scoped>
    .skill-container {
        display: flex;
        flex-direction: column;
        padding: 5px 10px;
    }

    .skill-container:hover {
        background-color: #111;
        cursor: pointer;
    }

    .skill-container .name-container {
        display: flex;
    }

    .skill-container .name-container .name {
        flex-grow: 1;
        font-size: 13px;
        font-weight: 600;
    }

    .skill-container .name-container .level {
        flex-basis: 50px;
        text-align: right;
        padding-right: 2px;
    }

    .skill-container .bar-container {
        border: 1px solid #3a3a3a;
        box-shadow: inset 0px 0px 4px #000;
        height: 10px;
        margin-top: 3px;
    }

    .skill-container .bar-container .bar {
        height: 100%;
        background-color: #424242;
        box-shadow: inset 0px 0px 2px #000;
    }
</style>