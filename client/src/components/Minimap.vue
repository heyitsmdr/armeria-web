<template>
    <div class="container">
        <div class="area-title" @click="handleClick">
            Some Area
        </div>
        <div class="map" ref="map"></div>
    </div>
</template>

<script>
export default {
    name: 'Minimap',
    data: () => {
        return {
            gridSize: 25,
            gridPadding: 7,
            mapHeight: 0,
            mapWidth: 0,
        }
    },
    methods: {
        handleClick() {
            this.renderMap(
                { x: 0, y: 0 },
                [
                    { x: 0, y: 0 },
                    { x: 1, y: 0 },
                    { x: 0, y: -1 },
                ]
            );
        },
        clearMap() {
            const map = this.$refs['map'];
            while (map.firstChild) {
                map.firstChild.remove();
            }
        },
        renderMap(currentLocation, map) {
            this.clearMap()
            map.forEach(m => {
                const div = document.createElement('div');
                div.style.height = this.gridSize + 'px';
                div.style.width = this.gridSize + 'px';
                div.style.position = 'absolute';
                div.style.top = -((m.y * this.gridSize) + this.gridPadding * m.y) + 'px';
                div.style.left = ((m.x * this.gridSize) + this.gridPadding * m.x) + 'px';
                div.style.backgroundColor = 'rgba(255,255,255,0.6)';
                div.style.border = '2px solid rgba(255,255,255,1)';
                div.style.borderRadius = '5px';
                div.className = 'room';

                this.$refs['map'].appendChild(div)
            })
        }
    },
    mounted() {
        this.mapHeight = this.$refs['map'].clientHeight;
        this.mapWidth = this.$refs['map'].clientWidth;
        console.log(this.mapHeight, this.mapWidth)
    }
}
</script>

<style lang="scss">
.container {
    height: 100%;
    display: flex;
    flex-direction: column;
}

.area-title {
    text-align: center;
    padding: 5px;
    background-color: #1b1b1b;
    border-bottom: 1px solid #313131;
    font-weight: 600;
    color: #fff;
}

.map {
    background-color: #131313;
    flex-grow: 1;
    position: relative;

    .room {
        transition: all .1s ease-in-out;

        &:hover {
            cursor:pointer;
            transform: scale(1.1)
        }
    }
}
</style>