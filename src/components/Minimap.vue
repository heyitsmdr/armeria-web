<template>
    <div class="root">
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
    import {mapState} from 'vuex';
    import * as PIXI from 'pixi.js';
    import {ease} from 'pixi-ease';

    export default {
        name: 'Minimap',
        data: () => {
            return {
                gridSize: 20,
                gridBorderSize: 2,
                gridPadding: 8,
                mapHeight: 0,
                mapWidth: 0,
                areaTitle: 'Unknown',
                app: null,
                mapContainer: null,
                backgroundTilingSprite: null
            }
        },
        computed: mapState(['minimapData', 'characterLocation', 'roomTitle', 'permissions']),
        watch: {
            minimapData(data) {
                this.areaTitle = data.name;
                this.renderMap();
                this.centerMapOnLocation(this.characterLocation);
            },
            characterLocation(newLocation, oldLocation) {
                if (newLocation.z !== oldLocation.z) {
                    this.renderMap();
                }
                this.centerMapOnLocation(newLocation);
            },
            permissions() {
                if (this.minimapData.rooms.length > 0) {
                    this.renderMap();
                }
            }
        },
        methods: {
            /**
             * Return hex color from an rgb string.
             * @param {String} rgb
             * @param {Number} darken=0
             */
            rgbToHex(rgb, darken = 0) {
                const a = rgb.split(',');
                let r = a[0];
                let g = a[1];
                let b = a[2];

                if (darken > 0) {
                    r -= darken;
                    g -= darken;
                    b -= darken;
                    if (r < 0) {
                        r = 0;
                    }
                    if (g < 0) {
                        g = 0;
                    }
                    if (b < 0) {
                        b = 0;
                    }
                }

                return PIXI.utils.rgb2hex([r / 255, g / 255, b / 255]);
            },

            /**
             * Remove all children from the map container.
             */
            clearMap() {
                this.mapContainer.removeChildren();
            },

            /**
             * Calculate localized offsets for a given room.
             * @param {Room} room
             */
            localRoomOffsets(room) {
                const gridSizeFull = this.gridSize + (this.gridBorderSize * 2);
                return {
                    x: ((room.x * gridSizeFull) + (this.gridPadding * room.x)),
                    y: -((room.y * gridSizeFull) + (this.gridPadding * room.y)),
                    size: gridSizeFull,
                }
            },

            /**
             * Draw a directional line between two rooms on a PIXI.Graphics.
             * @param {PIXI.Graphics} lineGraphics
             * @param {Room} srcRoom
             * @param {Room} targetRoom
             * @param {String} direction
             */
            drawRoomLine(lineGraphics, srcRoom, targetRoom, direction) {
                const srcOffsets = this.localRoomOffsets(srcRoom);
                const targetOffsets = this.localRoomOffsets(targetRoom);
                const lineWidth = 2;

                let startX, startY, endX, endY;
                switch (direction) {
                    case 'north':
                        startX = srcOffsets.x + (srcOffsets.size / 2);
                        startY = srcOffsets.y;
                        endX = targetOffsets.x + (targetOffsets.size / 2);
                        endY = targetOffsets.y + targetOffsets.size;
                        break;
                    case 'south':
                        startX = srcOffsets.x + (srcOffsets.size / 2);
                        startY = srcOffsets.y + srcOffsets.size;
                        endX = targetOffsets.x + (targetOffsets.size / 2);
                        endY = targetOffsets.y;
                        break;
                    case 'east':
                        startX = srcOffsets.x + srcOffsets.size;
                        startY = srcOffsets.y + (srcOffsets.size / 2);
                        endX = targetOffsets.x;
                        endY = targetOffsets.y + (targetOffsets.size / 2);
                        break;
                    case 'west':
                        startX = srcOffsets.x;
                        startY = srcOffsets.y + (srcOffsets.size / 2);
                        endX =  targetOffsets.x + targetOffsets.size;
                        endY = targetOffsets.y + (targetOffsets.size / 2);
                        break;
                }

                // If two rooms are the same color, draw their connecting lines as that color.
                let lineColor = this.rgbToHex('190,190,190');
                if (srcRoom.color === targetRoom.color) {
                    lineColor = this.rgbToHex(srcRoom.color);
                }

                lineGraphics
                    .lineStyle(lineWidth, lineColor)
                    .moveTo(startX, startY)
                    .lineTo(endX, endY);
            },

            /**
             * Render a directional indicator on the map container.
             * @param {Number} x
             * @param {Number} y
             * @param {String} direction
             */
            drawExternalRoomConnector(x, y, direction) {
                let s = PIXI.Sprite.from('./gfx/areaTransition.png');
                s.x = x + 12;
                s.y = y + 12;
                s.anchor.set(0.5);

                if (direction === 'east') {
                    s.rotation = 90 * (Math.PI / 180);
                }
                if (direction === 'south') {
                    s.rotation = 180 * (Math.PI / 180);
                }
                if (direction === 'west') {
                    s.rotation = -90 * (Math.PI / 180);
                }
                if (direction === 'up') {
                    s.x += 5;
                    s.y += 8;
                }
                if (direction === 'down') {
                    s.rotation = 180 * (Math.PI / 180);
                    s.x -= 5;
                    s.y -= 8;
                }
                this.mapContainer.addChild(s);
            },

            /**
             * Returns the matching room within the minimap rooms of the Vuex store.
             * @param {String} roomDirString
             * @returns {*}
             */
            roomAt(roomDirString) {
                const sections = roomDirString.split(",");

                for(let i = 0; i < this.minimapData.rooms.length; i++) {
                    const sameX = this.minimapData.rooms[i].x.toString() === sections[1];
                    const sameY = this.minimapData.rooms[i].y.toString() === sections[2];
                    const sameZ = this.minimapData.rooms[i].z.toString() === sections[3];

                    if (sameX && sameY && sameZ) {
                        return this.minimapData.rooms[i];
                    }
                }

                return null
            },

            /**
             * Renders the minimap.
             */
            renderMap() {
                this.clearMap();

                const lineGraphics = new PIXI.Graphics();
                this.mapContainer.addChild(lineGraphics);

                const filteredRooms = this.minimapData.rooms.filter(r => r.z === this.characterLocation.z);
                filteredRooms.forEach(room => {
                    let file;
                    switch (room.type) {
                        case 'track':
                            file = './gfx/trackTile.png';
                            break;
                        case 'bank':
                            file = './gfx/bankTile.png';
                            break;
                        case 'armor':
                            file = './gfx/armorTile.png';
                            break;
                        case 'sword':
                            file = './gfx/swordTile.png';
                            break;
                        case 'home':
                            file = './gfx/homeTile.png';
                            break;
                        case 'wand':
                            file = './gfx/wandTile.png';
                            break;
                        default:
                            file = './gfx/baseTile.png';
                    }

                    let sprite = PIXI.Sprite.from(file);
                    sprite.x = this.localRoomOffsets(room).x;
                    sprite.y = this.localRoomOffsets(room).y;
                    sprite.interactive = true;
                    sprite.buttonMode = true;
                    if (this.$store.state.permissions.indexOf('CAN_BUILD') >= 0) {
                        sprite.on('pointerdown', (e) => this.handleRoomClick(e, room));
                        sprite.on('pointerover', () => sprite.tint = this.rgbToHex('255,255,0'));
                        sprite.on('pointerout', () => sprite.tint = this.rgbToHex(room.color));
                    }
                    sprite.tint = this.rgbToHex(room.color);
                    this.mapContainer.addChild(sprite);

                    const directions = ['north', 'south', 'east', 'west', 'up', 'down'];
                    directions.forEach(dir => {
                        if (room[dir].length > 0) {
                            if (dir === 'up' || dir === 'down' || room[dir].split(',')[0] !== this.areaTitle) {
                                this.drawExternalRoomConnector(sprite.x, sprite.y, dir);
                            } else {
                                this.drawRoomLine(lineGraphics, room, this.roomAt(room[dir]), dir);
                            }
                        }
                    });
                });
            },

            /**
             * Centers the map using easing on the character's current location.
             * @param location
             */
            centerMapOnLocation: function (location) {
                const gridSizeFull = this.gridSize + (this.gridBorderSize * 2)
                var x = (this.app.screen.width / 2) + -((location.x * gridSizeFull) + (this.gridPadding * location.x)) - (gridSizeFull / 2);
                var y = (this.app.screen.height / 2) + ((location.y * gridSizeFull) + (this.gridPadding * location.y)) - (gridSizeFull / 2);
                ease.add(this.mapContainer, {x: x, y: y}, {duration: 100, repeat: false, reverse: false});

                // Tint the background.
                this.backgroundTilingSprite.tint = this.rgbToHex(this.roomColor(location), -100);

            },

            /**
             * Returns the room color for a particular location.
             * @param location
             */
            roomColor: function (location) {
                for (let roomIdx = 0; roomIdx < this.minimapData.rooms.length; roomIdx++) {
                    const room = this.minimapData.rooms[roomIdx];
                    if (room.x === location.x && room.y === location.y && room.z === location.z) {
                        return room.color;
                    }
                }

                return '0,0,0';
            },

            /**
             * Handles the click event for the area name container.
             * @param {MouseEvent} e
             */
            handleAreaClick: function (e) {
                if (e.shiftKey && this.$store.state.permissions.indexOf('CAN_BUILD') >= 0) {
                    this.$socket.sendObj({type: 'command', payload: '/area edit'});
                }
            },

            /**
             * Handles the click event for a minimap room.
             * @param {PIXI.InteractionEvent} evt
             * @param {Room} room
             */
            handleRoomClick: function (evt, room) {
                if (this.$store.state.permissions.indexOf('CAN_BUILD') >= 0) {
                    if (evt.data.originalEvent.shiftKey) {
                        this.$socket.sendObj({
                          type: 'command',
                          payload: `/room edit ${room.x},${room.y},${room.z}`
                        });
                    } else {
                        this.$socket.sendObj({
                          type: 'command',
                          payload: `/teleport ${this.areaTitle},${room.x},${room.y},${room.z}`
                        });
                    }

                }
            }
        },

        mounted() {
            const mapCanvas = document.getElementById('map-canvas');
            this.app = new PIXI.Application({width: 250, height: 205, view: mapCanvas, antialias: true});

            this.backgroundTilingSprite = new PIXI.TilingSprite(
                PIXI.Texture.from('gfx/minimap-bg.png'),
                this.app.screen.width,
                this.app.screen.height
            );

            const backgroundContainer = new PIXI.Container();
            backgroundContainer.addChild(this.backgroundTilingSprite);

            const filter = new PIXI.filters.ColorMatrixFilter();
            backgroundContainer.filters = [filter];
            filter.contrast(1, true);
            this.app.stage.addChild(backgroundContainer);

            this.mapContainer = new PIXI.Container();
            this.app.stage.addChild(this.mapContainer);

            const pos = this.$refs['position'];

            // set position to half height/width with an offset for border size on position marker
            pos.style.top = ((this.app.screen.height / 2) - 2) + 'px';
            pos.style.left = ((this.app.screen.width / 2) - 2) + 'px';
        }
    }
</script>

<style scoped lang="scss">
    @import "@/styles/common";

    .root {
        height: 100%;
        display: flex;
        flex-direction: column;
        box-sizing: border-box;
        /*border: $defaultBorder;*/
        @include defaultBorderImage;
    }

    .area-title {
        text-align: center;
        padding: 5px;
        background-color: $bg-color;
        border-bottom: $defaultBorder;
        color: #fff;
        flex-shrink: 1;
    }

    .area-title .map-name {
        font-weight: 600;
        font-size: 16px;
        cursor: pointer;
    }

    .area-title .room-name {
        font-size: 12px;
        color: #bdbdbd;
    }

    .map {
        background-color: $bg-color;
        flex-basis: 205px;
        position: relative;
        overflow: hidden;
        transition: all .3s ease-in-out;
        box-shadow: inset 0px 0px 10px 1px #000;
    }

    .map .position {
        position: absolute;
        height: 0px;
        width: 0px;
        background-color: #ff0;
        border: 2px solid #FF0;
    }

    .map .floor {
        position: absolute;
        top: 0px;
        left: 0px;
        transition: all .1s ease-in-out;
    }
</style>
