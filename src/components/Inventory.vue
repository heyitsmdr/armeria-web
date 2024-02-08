<template>
  <div class="root">
    <div class="item-container">
      <Item
        v-for="item in items"
        :key="item.slot"
        :uuid="item.uuid"
        :name="item.name"
        :slotNum="item.slot"
        :equipSlot="item.equipSlot"
        :pictureKey="item.picture"
        :color="item.color"
      />
    </div>
    <div class="currency-container"><b>Money:</b> {{ formattedMoney }}</div>
  </div>
</template>

<script setup>
import { ref, computed, onBeforeMount } from "vue";
import { useStore } from "vuex";
import Item from "@/components/Item.vue";

const formatter = ref(null);

const store = useStore();
const inventory = computed(() => store.state.inventory);
const money = computed(() => store.state.money);

const formattedMoney = computed(() => {
  return formatter.value.format(parseFloat(money.value));
});
const items = computed(() => {
  let itemDef = {};
  inventory.value.forEach((item) => {
    itemDef[item.slot] = item;
  });

  let items = [];
  for (let i = 0; i < 35; i++) {
    if (itemDef[i]) {
      items.push(itemDef[i]);
    } else {
      items.push({ slot: i });
    }
  }
  return items;
});

onBeforeMount(() => {
  formatter.value = new Intl.NumberFormat("en-US", {
    style: "currency",
    currency: "USD",
  });
});
</script>

<style scoped lang="scss">
@import "@/styles/common";

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
