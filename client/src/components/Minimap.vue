<template>
    <div class="container">
        <div class="area-title">
            <div class="map-name" @click="handleAreaClick">{{ areaTitle }}</div>
            <div class="room-name">{{ roomTitle }}</div>
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
            gridSize: 22,
            gridBorderSize: 2,
            gridPadding: 6,
            mapHeight: 0,
            mapWidth: 0,
            areaTitle: 'Unknown'
        }
    },
    watch: {
        minimapData: function(newMinimapData) {
            this.areaTitle = newMinimapData.name;
            this.renderMap(newMinimapData.rooms, this.characterLocation.z);
            this.centerMapOnLocation(this.characterLocation, this.characterLocation);
        },
        characterLocation: function(newLocation, oldLocation) {
            if (newLocation.z !== oldLocation.z) {
                this.renderMap(this.minimapData.rooms, newLocation.z);
            }
            this.centerMapOnLocation(newLocation, oldLocation);
        }
    },
    computed: mapState(['minimapData', 'characterLocation', 'roomTitle']),
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

            const filteredRooms = rooms.filter(r => r.z === zIndex);

            filteredRooms.forEach(room => {
                const div = document.createElement('div');
                div.style.height = this.gridSize + 'px';
                div.style.width = this.gridSize + 'px';
                div.style.position = 'absolute';
                div.style.top = -((room.y * gridSizeFull) + (this.gridPadding * room.y)) + 'px';
                div.style.left = ((room.x * gridSizeFull) + (this.gridPadding * room.x)) + 'px';
                div.style.backgroundColor = `rgba(${room.color},0.6)`;
                div.style.border = `${this.gridBorderSize}px solid rgba(${room.color},0.8)`;
                div.setAttribute('x', room.x);
                div.setAttribute('y', room.y);
                div.setAttribute('z', room.z);
                div.addEventListener('mouseover', this.onRoomHover);
                div.addEventListener('click', this.handleRoomClick);
                div.className = 'room';

                if (room.type === 'track') {
                    div.style.opacity = '0.3';
                    div.style.borderRadius = '20px';
                }

                this.$refs['floor'].appendChild(div)
            })
        },

        centerMapOnLocation(location, oldLocation) {
            const floor = this.$refs['floor'];
            const halfMapWidth = this.mapWidth / 2;
            const halfMapHeight = this.mapHeight / 2;
            const gridSizeFull = this.gridSize + (this.gridBorderSize * 2);

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

            for(let i = 0; i < this.minimapData.rooms.length; i++) {
                const x = this.minimapData.rooms[i].x;
                const y = this.minimapData.rooms[i].y;
                const z = this.minimapData.rooms[i].z;
                if (x === location.x && y === location.y && z === location.z) {
                    this.$refs['map'].style.backgroundColor = `rgba(${this.minimapData.rooms[i].color},0.05)`;
                    break;
                }
            }
        },

        handleAreaClick: function(e) {
            if (e.shiftKey) {
                this.$socket.sendObj({type: 'command', payload: '/area edit'});
            }
        },

        handleRoomClick: function(e) {
            if (e.shiftKey) {
                var room = e.srcElement;
                var coords = room.getAttribute('x') +
                        ',' + room.getAttribute('y') +
                        ',' + room.getAttribute('z');
                this.$socket.sendObj({type: 'command', payload: '/room edit ' + coords});
            }
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
         transform: scale(1.3);
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
    color: #fff;
    flex-shrink: 1;

    .map-name {
        font-weight: 600;
        font-size: 16px;
        cursor: pointer;
    }

    .room-name {
        font-size: 12px;
    }
}

.map {
    background-color: #0c0c0c;
    border-bottom: 1px solid #313131;
    flex-grow: 1;
    min-height: 205px;
    position: relative;
    overflow: hidden;
    transition: all .3s ease-in-out;
    box-shadow: inset 0px 0px 10px 1px #000;

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
