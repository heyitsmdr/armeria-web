<template>
    <div class="container">
        <div class="area-title">
            {{ areaTitle }}
        </div>
        <div class="map" ref="map">
            <div class="floor" ref="floor">

            </div>
            <div class="position" ref="position"></div>
        </div>
    </div>
</template>

<script>
import { mapState } from 'vuex';

export default {
    name: 'Minimap',
    data: () => {
        return {
            gridSize: 24,
            gridBorderSize: 2,
            gridPadding: 8,
            mapHeight: 0,
            mapWidth: 0,
            areaTitle: '-',
        }
    },
    watch: {
        minimapData: function(newMinimapData) {
            this.renderMap(newMinimapData.rooms, 0);
            this.areaTitle = newMinimapData.name;
        },
        characterLocation: function(newLocation, oldLocation) {
            this.centerMapOnLocation(newLocation, oldLocation);
        }
    },
    computed: mapState(['minimapData', 'characterLocation']),
    methods: {
        clearMap() {
            const map = this.$refs['floor'];
            while (map.firstChild) {
                map.firstChild.remove();
            }
        },
        renderMap(rooms, zIndex) {
            const gridSizeFull = this.gridSize + (this.gridBorderSize * 2)

            this.clearMap();

            const filteredRooms = rooms.filter(r => r.z == zIndex);

            filteredRooms.forEach(room => {
                const div = document.createElement('div');
                div.style.height = this.gridSize + 'px';
                div.style.width = this.gridSize + 'px';
                div.style.position = 'absolute';
                div.style.top = -((room.y * gridSizeFull) + (this.gridPadding * room.y)) + 'px';
                div.style.left = ((room.x * gridSizeFull) + (this.gridPadding * room.x)) + 'px';
                div.style.backgroundColor = 'rgba(255,255,255,0.6)';
                div.style.border = `${this.gridBorderSize}px solid rgba(255,255,255,1)`;
                div.setAttribute('x', room.x);
                div.setAttribute('y', room.y);
                div.setAttribute('z', room.z);
                div.addEventListener('mouseover', this.onRoomHover);
                div.className = 'room';

                this.$refs['floor'].appendChild(div)
            })
        },
        centerMapOnLocation(location, oldLocation) {
            const floor = this.$refs['floor'];
            const halfMapWidth = this.mapWidth / 2;
            const halfMapHeight = this.mapHeight / 2;
            const gridSizeFull = this.gridSize + (this.gridBorderSize * 2)

            floor.style.left = (halfMapWidth - (gridSizeFull / 2) - (gridSizeFull * location.x) - (this.gridPadding * location.x)) + 'px';
            floor.style.top = (halfMapHeight - (gridSizeFull / 2) - (gridSizeFull * -location.y) - (this.gridPadding * -location.y)) + 'px';

            const oldLocDiv = document.querySelector(`.room[x="${oldLocation.x}"][y="${oldLocation.y}"]`)
            if  (oldLocDiv) {
                oldLocDiv.classList.remove('current-location');
            }
            const newLocDiv = document.querySelector(`.room[x="${location.x}"][y="${location.y}"]`)
            if (newLocDiv) {
                newLocDiv.classList.add('current-location');
            }
            this.location = location;
        },
        onRoomHover(event) {
            if (!event.target.classList.contains('current-location')) {
                this.$playSound('bloop.wav');
            }
        }
    },
    mounted() {
        const map = this.$refs['map'];
        const pos = this.$refs['position'];
        this.mapHeight = map.clientHeight;
        this.mapWidth = map.clientWidth;
        // set position to half height/width with an offset for border size on position marker
        pos.style.top = ((this.mapHeight / 2) - 2) + 'px';
        pos.style.left = ((this.mapWidth / 2) - 2) + 'px';
    }
}
</script>

<style lang="scss">
.map .floor .room {
    transition: all .1s ease-in-out;

    &:hover {
         cursor: pointer;
         transform: scale(1.1);
     }

    &.current-location {
         border-color: #ff0 !important;
         transform: scale(1.2);
    }
}
</style>

<style lang="scss" scoped>
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
    font-size: 16px;
    color: #fff;
}

.map {
    background-color: #0c0c0c;
    border-bottom: 1px solid #313131;
    flex-grow: 1;
    position: relative;
    overflow: hidden;

    .position {
        position: absolute;
        height: 0px;
        width: 0px;
        background-color: #ff0;
        border: 2px solid #FF0;
    }

    .floor {
        position: absolute;
        top: 0px;
        left: 0px;
        transition: all .1s ease-in-out;
    }
}
</style>