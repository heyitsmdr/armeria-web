<template>
    <div class="skills-container">
        <div class="skill-level-container">
            <div class="text">Skill Level</div>
            <div class="level">100</div>
        </div>
        <div class="skill-items-container">
            <Skill name="Woodworking" level="15" percent="20" />
            <Skill name="Climbing" level="20" percent="50" />
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

<style scoped>
    .skills-container {
        display: flex;
        flex-direction: column;
    }

    .skills-container .skill-level-container {
        flex-basis: 40px;
        display: flex;
        box-shadow: 0px 0px 7px #000;
    }

    .skills-container .skill-level-container .text {
        flex-grow: 1;
        display: flex;
        align-items: center;
        font-size: 18px;
        padding-left: 10px;
        background-color: #222;
    }

    .skills-container .skill-level-container .level {
        flex-basis: 60px;
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 20px;
        font-weight: 600;
        background-color: #111;
        color: #03A9F4;
    }

    .skills-container .skill-items-container {
        flex-grow: 1;
    }
</style>