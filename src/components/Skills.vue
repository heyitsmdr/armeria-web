<template>
    <div class="root">
        <div class="skill-level-container">
            <div class="text">Skill Level</div>
            <div class="level">0</div>
        </div>
        <div class="skill-items-container">
            <!--<div class="no-skills">No Skills</div>-->
            <Skill name="Test Skill" level="7"/>
        </div>
    </div>
</template>

<script>
    import Skill from '@/components/Skill';

    export default {
        name: 'Skills',
        components: { Skill },
    }
</script>

<style scoped lang="scss">
    @import "@/styles/common";

    .root {
        display: flex;
        flex-direction: column;
        height: 100%;
        box-sizing: border-box;
        border: $defaultBorder;
    }

    .root .skill-level-container {
        flex-basis: 40px;
        display: flex;
    }

    .root .skill-level-container .text {
        flex-grow: 1;
        display: flex;
        align-items: center;
        font-size: 16px;
        padding-left: 10px;
        background-color: $defaultBackgroundColor;
    }

    .root .skill-level-container .level {
        flex-basis: 60px;
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 20px;
        font-weight: 600;
        background-color: $defaultBackgroundColor;
        color: $defaultTextColor;
    }

    .root .skill-items-container {
        flex-grow: 1;
    }

    .root .skill-items-container .no-skills {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        color: $defaultTextColor;
    }
</style>