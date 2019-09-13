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
import { Ease, ease } from 'pixi-ease';

export default {
    name: 'Minimap',
    data: () => {
        return {
            gridSize: 22,
            gridBorderSize: 2,
            gridPadding: 12,
            mapHeight: 0,
            mapWidth: 0,
            areaTitle: 'Unknown'
        }
    },
    watch: {
        minimapData: function(newMinimapData) {
            this.areaTitle = newMinimapData.name;
            
            this.renderMap(newMinimapData.rooms, this.characterLocation, this.$socket, this.$store);
            this.centerMapOnLocation(this.characterLocation);
        },
        characterLocation: function(newLocation, oldLocation) {
            if (newLocation.z !== oldLocation.z) {
                this.renderMap(this.minimapData.rooms, this.characterLocation, this.$socket, this.$store);
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

        toHex(c) {
            var hex = c.toString(16);
            return hex.length == 1 ? "0" + hex : hex;
        },

        clearMap() {
            this.mapContainer.removeChildren();
        },

        renderMap(rooms, loc, socket, store) {
            this.app.stage.addChild(this.mapContainer);

            const gridSizeFull = this.gridSize + (this.gridBorderSize * 2)
            
            this.clearMap();

            const filteredRooms = rooms.filter(r => r.z === loc.z);

            filteredRooms.forEach(room => {
                let sprite = PIXI.Sprite.from('./gfx/baseTile.png');
                sprite.anchor.set(0.5);
                sprite.x = ((room.x * gridSizeFull) + (this.gridPadding * room.x));
                sprite.y = -((room.y * gridSizeFull) + (this.gridPadding * room.y));
                sprite.interactive = true;
                sprite.buttonMode = true;
                sprite.on('pointerdown', handleRoomClick);
                sprite.on('pointerover', function(){sprite.scale.set(1.2, 1.2)});
                sprite.on('pointerout', function(){sprite.scale.set(1.0,1.0)});
                this.mapContainer.addChild(sprite);

                sprite.tint = this.rgbToHex(room.color);

                function handleRoomClick(){
                    if (store.state.permissions.indexOf('CAN_BUILD') >= 0) {
                        socket.sendObj({type: 'command', payload: '/room edit ' + room.x + ',' + room.y + ',' + room.z});
                    }
                }
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
        }
    },
    mounted() {
        const mapCanvas = document.getElementById("map-canvas");
        this.app = new PIXI.Application({width: 250, height: 206, view: mapCanvas});
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
