<template>
    <div
        class="tooltip"
        :class="{ visible: itemTooltipVisible }"
        :style="{ borderColor: rarityColor }"
        ref="tooltip"
        v-html="htmlData"
    ></div>
</template>

<script>
    import {mapState} from 'vuex';

    export default {
        name: 'ItemTooltip',
        computed: mapState(['itemTooltipVisible', 'itemTooltipUUID', 'itemTooltipCache', 'itemTooltipMouseCoords']),
        data: () => {
            return {
                itemUUID: '',
                htmlData: '',
                rarityColor: '',
            }
        },
        watch: {
            itemTooltipUUID: function(uuid) {
                if (uuid.length === 0) {
                    this.htmlData = '';
                    this.itemUUID = uuid;
                    this.rarityColor = '';
                } else if (this.itemUUID !== uuid) {
                    this.itemUUID = uuid;

                    const renderSuccess = this.renderHTML();
                    if (!renderSuccess) {
                        this.$socket.sendObj({
                            type: 'itemTooltipHTML',
                            payload: this.itemUUID
                        });
                    }
                }
            },

            itemTooltipCache: function() {
                this.renderHTML();
            },

            itemTooltipVisible: function(visible) {
                if (visible) {
                    const tt = this.$refs["tooltip"];
                    tt.style.top = '-500px';
                    tt.style.left = '-500px';
                }
            },
            itemTooltipMouseCoords: function(coords) {
                const tt = this.$refs["tooltip"];

                if (tt.clientHeight === 0) {
                    return;
                }

                const xOffset = 50;
                const yOffset = 50;

                let ttTop = coords.y - tt.clientHeight - yOffset;
                let ttLeft = coords.x - xOffset;

                if ((ttLeft + tt.clientWidth + 20) > window.innerWidth) {
                    ttLeft = window.innerWidth - tt.clientWidth - 20;
                } else if (ttLeft < 10) {
                    ttLeft = 10;
                }


                tt.style.top = ttTop + 'px';
                tt.style.left = ttLeft + 'px';
            }
        },
        methods: {
            renderHTML: function() {
                const cachedItem = this.$store.getters.itemTooltipCache(this.itemUUID);

                if (cachedItem) {
                    this.htmlData = cachedItem.html;
                    this.rarityColor = cachedItem.rarity;

                    return true
                }

                return false
            }
        }
    }
</script>

<style scoped>
    .tooltip {
        display: none;
        position: absolute;
        max-width: 400px;
        min-width: 150px;
        z-index: 999;
        background-color: #111;
        border: 2px solid #ccc;
        border-radius: 5px;
        padding: 5px;
        box-shadow: 0px 0px 10px #000;
    }

    .tooltip.visible {
        display: block;
    }
</style>