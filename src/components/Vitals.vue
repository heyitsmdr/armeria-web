<template>
    <div class="targets-container">
        <div class="bar-container">
            <div class="health" ref="health-bar"></div>
            <div class="text">
                <div class="visible">Health - 100%</div>
                <div class="extended">1000 / 1000</div>
            </div>
            <div class="border-overlay" @mouseover="handleMouseOver" @mouseout="handleMouseOut"></div>
        </div>
        <div class="bar-container">
            <div class="magic" ref="magic-bar"></div>
            <div class="text">
                <div class="visible">Magic - 100%</div>
                <div class="extended">1000 / 1000</div>
            </div>
            <div class="border-overlay" @mouseover="handleMouseOver" @mouseout="handleMouseOut"></div>
        </div>
        <div class="bar-container">
            <div class="stamina" ref="stamina-bar"></div>
            <div class="text">
                <div class="visible">Stamina - 100%</div>
                <div class="extended">1000 / 1000</div>
            </div>
            <div class="border-overlay" @mouseover="handleMouseOver" @mouseout="handleMouseOut"></div>
        </div>
    </div>
</template>

<script>
    export default {
        name: 'Vitals',
        methods: {
            handleMouseOver: function(e) {
                const textSibling = e.target.parentNode.querySelector('.text');
                textSibling.classList.add('show-extended');
            },

            handleMouseOut: function(e) {
                const textSibling = e.target.parentNode.querySelector('.text');
                textSibling.classList.remove('show-extended');
            },
        },
    }
</script>

<style scoped>
    .targets-container {
        display: flex;
        height: 29px;
        background-color:#0c0c0c;
    }

    .targets-container .bar-container {
        flex-basis: 33.333%;
        background-color: rgba(0,0,0,0);
        z-index: 2;
        position: relative;
        overflow: hidden;
    }

    .targets-container .bar-container .border-overlay {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: rgba(0,0,0,0);
        box-shadow: inset 0px 0px 5px 4px #000;
    }

    .targets-container .bar-container .border-overlay:hover {
        cursor: pointer;
    }

    .targets-container .bar-container .health {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: #e91e63;
        border-right: 2px solid #e0628c;
    }

    .targets-container .bar-container .magic {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: #03a9f4;
        border-right: 2px solid #6dc5ec;
    }

    .targets-container .bar-container .stamina {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        background-color: #ff5722;
        border-right: 2px solid #e2a895;
    }

    .targets-container .bar-container .text {
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        color: #fff;
        display: flex;
        justify-content: center;
        align-items: center;
        font-size: 12px;
        text-shadow: 0px 0px 4px #000;
    }

    .targets-container .bar-container .text .visible {
        position: absolute;
        transition: all 0.1s ease-in-out;
        top: 6px;
    }

    .targets-container .bar-container .text .extended {
        position: absolute;
        transition: all 0.1s ease-in-out;
        top: -50px;
    }

    .targets-container .bar-container .text.show-extended .extended {
        top: 6px;
    }

    .targets-container .bar-container .text.show-extended .visible {
        top: 50px;
    }
</style>