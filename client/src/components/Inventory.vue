<template>
    <div class="inventory">
        <div class="item-container">
            <Item v-for="item in items" />
        </div>
    </div>
</template>

<script>
    import Item from '@/components/Item';

    export default {
        name: 'Inventory',
        components: {
            Item
        },
        computed: {
            items: function() {
                let items = [];
                for(let i = 0; i < 35; i++) {
                    items.push('');
                }
                return items;
            }
        }
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
        padding: 10px;
    }
</style>