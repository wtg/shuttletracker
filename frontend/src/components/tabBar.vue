<template>
  <ul>
    <router-link tag="li" to="/">
      <span>
        <span class="fas fa-location-arrow"></span>
        Map
      </span>
    </router-link>
    <router-link tag="li" to="/schedules">
      <span>
        <span class="fas fa-list-ul"></span>
        Schedules
      </span>
    </router-link>
    <!-- <router-link v-if="etasEnabled" tag="li" to="/etas">
      ETAs
    </router-link> -->
    <router-link tag="li" to="/settings">
      <span>
        <span class="fas fa-cog"></span>
        Settings
      </span>
    </router-link>
  </ul>
</template>


<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  computed: {
    etasEnabled(): boolean {
      return this.$store.state.settings.etasEnabled;
    },
  },
});
</script>

<style lang="scss" scoped>
@import "@/assets/vars.scss";
ul {
  display: flex;
  height: 40px;
  position: fixed;
  bottom: 0;
  display: flex;
  padding: 0;
  margin: 0;
  justify-content: center;
  align-items: center;
  width: 100%;
  border-top: 0.5px solid #eee;
  box-shadow: 0 3px 8px 0 #ddd;
  font-size: 13px;
  user-select: none;
  background: white;
}
li {
  cursor: pointer;
  width: auto;
  height: 100%;
  padding: 5px 15px;
  margin: 0 5px;
  border-top: 1px solid rgba(0, 0, 0, 0);
  position: relative;
  top: -0.5px;
  display: flex;
  align-items: center;
}
.router-link-exact-active {
  border-top: 1px solid $primary;
  color: $primary;
}
</style>
