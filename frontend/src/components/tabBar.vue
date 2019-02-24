<template>
  <ul>
    <router-link tag="li" to="/">Map</router-link>
    <router-link tag="li" to="/schedules">Schedules</router-link>
    <router-link tag="li" to="/settings">Settings</router-link>
  </ul>
</template>


<script lang="ts">
import Vue from 'vue';

export default Vue.extend({});
</script>

<style lang="scss" scoped>
ul {
  display: flex;
  height: 40px;
  position: fixed;
  bottom: 0;
  display: flex;
  padding: 0;
  margin: 0;
  justify-content: space-evenly;
  align-items: center;
  width: 100%;
  border-top: 1px solid #eee;
  box-shadow: 0 3px 8px 0 #ddd;
  font-size: 13px;
  user-select: none;
  background: white;
}
li {
  cursor: pointer;
}
.router-link-exact-active {
  font-weight: 600;
}
</style>
