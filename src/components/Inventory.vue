<template>
    <div class="root">
        <div class="item-container">
            <Item
                    v-for="item in items"
                    :key="item.slot"
                    :uuid="item.uuid"
                    :name="item.name"
                    :slotNum="item.slot"
                    :pictureKey="item.picture"
                    :color="item.color"
            />
        </div>
        <div class="currency-container">
            <b>Money:</b> {{ formattedMoney }}
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

<style scoped lang="scss">
    @import "~@/styles/common";

    .root {
        height: 100%;
        box-sizing: border-box;
        /*border: $defaultBorder;*/
        @include defaultBorderImage;
    }

    .item-container {
        display: flex;
        flex-wrap: wrap;
        padding: 10px;
        justify-content: space-evenly;
    }

    .currency-container {
        color: #ffc107;
        font-size: 13px;
    }
</style>