<template>
    <div class="skills-container">
        <div class="skill-level-container">
            <div class="text">Skill Level</div>
            <div class="level">0</div>
        </div>
        <div class="skill-items-container">
            <div class="no-skills">No Skills</div>
        </div>
    </div>
</template>

<script>
    //import Skill from '@/components/Skill';

    export default {
        name: 'Skills',
        //components: { Skill },
    }
</script>

<style scoped>
    .skills-container {
        display: flex;
        flex-direction: column;
        height: 100%;
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
        font-size: 16px;
        padding-left: 10px;
        background-color: #1b1b1b;
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

    .skills-container .skill-items-container .no-skills {
        display: flex;
        justify-content: center;
        align-items: center;
        height: 100%;
        color: #666;
    }
</style>