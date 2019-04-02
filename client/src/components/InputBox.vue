<template>
  <div class="container">
    <input class="input-box" type="text" />
  </div>
</template>

<script>
export default {
  name: 'InputBox'
}
</script>

<style lang="scss" scoped>
.container {
  padding-left: 5px;
}

.input-box {
  background-color: #000;
  border: 0;
  height: 35px;
  width: 100%;
  color: #fff;
  font-family: 'Montserrat', sans-serif;
  font-weight: 500;
  font-size: 13px;

  &:focus{
    outline: none;
  }
}
</style>