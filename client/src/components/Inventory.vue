<template>
    <div class="inventory">
        <div class="item-container">
            <Item
                    v-for="item in items"
                    :key="item.slot"
                    :uuid="item.uuid"
                    :slotNum="item.slot"
                    :pictureKey="item.picture"
                    :color="item.color"
                    :tooltipData="item.tooltip"
            />
        </div>
        <div class="currency-container">
            <b>Money:</b> $0
        </div>
    </div>
</template>

<script>
    import {mapState} from 'vuex';
    import Item from '@/components/Item';

    export default {
        name: 'Inventory',
        components: {
            Item
        },
        computed: {
            ...mapState(['inventory']),
            items: function () {
                let itemDef = {};
                this.inventory.forEach(item => {
                    itemDef[item.slot] = item
                });

                let items = [];
                for (let i = 0; i < 35; i++) {
                    if (itemDef[i]) {
                        items.push(itemDef[i]);
                    } else {
                        items.push({slot: i});
                    }
                }
                return items;
            }
        },
    }
</script>

<style scoped>
    .inventory {
        background-color: #131313;
        height: 100%;
    }

    .item-container {
        display: flex;
        flex-wrap: wrap;
        padding: 10px 10px 0px 10px;
    }

    .currency-container {
        margin-left: 10px;
        margin-right: 14px;
        border: 1px solid #bf9109;
        color: #ffc107;
        padding: 4px;
        box-shadow: inset 0px 0px 5px #cc9904;
        background-color: #000;
        font-size: 13px;
    }
</style>