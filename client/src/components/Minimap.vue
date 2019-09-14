<template>
    <div class="container">
        <div class="area-title">
            <div class="map-name" @click="handleAreaClick">{{ areaTitle }}</div>
            <div class="room-name">{{ roomTitle }}</div>
        </div>
        <div class="map" ref="map">
            <canvas id="map-canvas"></canvas>
            <div class="floor" ref="floor">

            </div>
            <div class="position" ref="position"></div>
        </div>
    </div>
</template>

<script>
import { mapState } from 'vuex';
import * as PIXI from 'pixi.js';
import { ease } from 'pixi-ease';

export default {
    name: 'Minimap',
    data: () => {
        return {
            gridSize: 22,
            gridBorderSize: 2,
            gridPadding: 2,
            mapHeight: 0,
            mapWidth: 0,
            areaTitle: 'Unknown',
            app: null,
            mapContainer: null
        }
    },
    watch: {
        minimapData: function(newMinimapData) {
            this.areaTitle = newMinimapData.name;
            this.renderMap(newMinimapData.rooms, this.characterLocation);
            this.centerMapOnLocation(this.characterLocation);
        },
        characterLocation: function(newLocation, oldLocation) {
            if (newLocation.z !== oldLocation.z) {
                this.renderMap(this.minimapData.rooms, this.characterLocation);
            }
            this.centerMapOnLocation(newLocation);
        }
    },
    computed: mapState(['minimapData', 'characterLocation', 'roomTitle']),
    methods: {
        rgbToHex(rgb) {
            var a = rgb.split(',');
            return PIXI.utils.rgb2hex([a[0]/255, a[1]/255, a[2]/255]);
        },

        clearMap() {
            this.mapContainer.removeChildren();
        },

        renderMap(rooms, loc) {
            this.app.stage.addChild(this.mapContainer);

            const gridSizeFull = this.gridSize + (this.gridBorderSize * 2)
            
            this.clearMap();

            const filteredRooms = rooms.filter(r => r.z === loc.z);

            filteredRooms.forEach(room => {
                let file;
                switch (room.type) {
                    case 'track':
                        file = './gfx/trackTile.png';
                        break;
                    default:
                        file = './gfx/baseTile.png';
                }
                let sprite = PIXI.Sprite.from(file);
                sprite.anchor.set(0.5);
                sprite.x = ((room.x * gridSizeFull) + (this.gridPadding * room.x));
                sprite.y = -((room.y * gridSizeFull) + (this.gridPadding * room.y));
                sprite.interactive = true;
                sprite.buttonMode = true;
                sprite.on('pointerdown', () => this.handleRoomClick(room));
                sprite.on('pointerover', () => sprite.scale.set(1.2, 1.2));
                sprite.on('pointerout', () => sprite.scale.set(1.0,1.0));
                this.mapContainer.addChild(sprite);
                sprite.tint = this.rgbToHex(room.color);

                let areaTransitions = [];
                if (room.north !== '') { areaTransitions.push('n'); }
                if (room.east !== '') { areaTransitions.push('e'); }
                if (room.south !== '') { areaTransitions.push('s'); }
                if (room.west !== '') { areaTransitions.push('w'); }
                if (room.up !== '') { areaTransitions.push('u'); }
                if (room.down !== '') { areaTransitions.push('d'); }

                areaTransitions.forEach(t => {
                    let s = PIXI.Sprite.from('./gfx/areaTransition.png');
                    s.x = sprite.x;
                    s.y = sprite.y;
                    s.anchor.set(0.5);
                    if (t === 'e') { s.rotation = 90 * (Math.PI/180); }
                    if (t === 's') { s.rotation = 180 * (Math.PI/180); }
                    if (t === 'w') { s.rotation = -90 * (Math.PI/180); }
                    if (t === 'u') {
                        s.x = Math.floor(sprite.x + gridSizeFull / 4) + 1;
                        s.y = Math.floor(sprite.y + gridSizeFull / 4);
                    }
                    if (t === 'd') {
                        s.rotation = 180 * (Math.PI/180);
                        s.x = Math.floor(sprite.x - gridSizeFull / 4);
                        s.y = Math.floor(sprite.y - gridSizeFull / 4) + 1;
                    }
                    this.mapContainer.addChild(s);
                })
            })
        },

        centerMapOnLocation: function(location) {
            const gridSizeFull = this.gridSize + (this.gridBorderSize * 2)
            var x = (this.app.screen.width / 2) + -((location.x * gridSizeFull) + (this.gridPadding * location.x));
            var y = (this.app.screen.height / 2) + ((location.y * gridSizeFull) + (this.gridPadding * location.y));
            ease.add(this.mapContainer, { x: x, y: y }, { duration: 100, repeat: false, reverse: false });//move(this.mapContainer, x, y);
        },

        handleAreaClick: function(e) {
            if (e.shiftKey && this.$store.state.permissions.indexOf('CAN_BUILD') >= 0) {
                this.$socket.sendObj({type: 'command', payload: '/area edit'});
            }
        },

        handleRoomClick: function(room) {
                    if (this.$store.state.permissions.indexOf('CAN_BUILD') >= 0) {
                        this.$socket.sendObj({type: 'command', payload: '/room edit ' + room.x + ',' + room.y + ',' + room.z});
                    }
                }
    },

    mounted() {
        const mapCanvas = document.getElementById('map-canvas');
        this.app = new PIXI.Application({width: 250, height: 206, view: mapCanvas, antialias: true});
        this.mapContainer = new PIXI.Container();

        const pos = this.$refs['position'];

        // set position to half height/width with an offset for border size on position marker
        pos.style.top = ((this.app.screen.height / 2) - 2) + 'px';
        pos.style.left = ((this.app.screen.width / 2) - 2) + 'px';
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
