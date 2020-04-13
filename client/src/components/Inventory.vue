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
            />
        </div>
        <div class="currency-container">
            <img src="gfx/iconGold.png" alt=""><p>{{ formattedMoney }}</p>
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
        data: function () {
            return {
                formatter: null,
            }
        },
        beforeMount: function () {
            this.formatter = new Intl.NumberFormat('en-US', {
                style: 'currency',
                currency: 'USD',
            });
        },
        computed: {
            ...mapState(['inventory', 'money']),
            formattedMoney: function() {
                return this.formatter.format(parseFloat(this.money));
            },
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
        background-color: #1b1b1b;
        height: 100%;
    }

    .item-container {
        display: flex;
        flex-wrap: wrap;
        padding: 5px 10px 0px 10px;
    }

    img {
            max-width: 100%;
            max-height: 100%;
            float:right;
        }

    .currency-container {
        margin-left: 10px;
        margin-right: 14px;
        color: #ffc107;
        background-color: #1b1b1b;
        font-size: 13px;
        height: 24px;
        margin-top: -12px;
    }

    p {
        height: 24px;
        display: flex;
        text-align: right;
        flex-direction: column;
        padding: 3px 3px 0 0;
    }
</style>