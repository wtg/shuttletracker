<template>
<div id="outer">
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
    <div class="version">
      <div class="demo version-section">
        <a href="https://github.com/wtg/shuttletracker" class="github-corner">
          <svg width="30" height="30" xmlns="http://www.w3.org/2000/svg" xmlns:svg="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 750 700" style="position: absolute; bottom: 0; border: 0; right: 0;">
            <g>
              <image stroke="null" xlink:href="../../public/githublogo.png" id="svg_2" height="477.99998" width="636.00002" y="2.00001" x="2"/>
            </g>
          </svg>
        </a>
      </div>
    </div>
  </ul>
</div>
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
  padding-right: 30px; //Account for github logo
  margin: 0;
  justify-content: center;
  align-items: center;
  width: 100%;
  // border-top: 0.5px solid var(--color-bg-less);
  box-shadow: 0 3px 8px 0 var(--color-bg-least);
  font-size: 13px;
  user-select: none;
  background: var(--color-bg-normal);
  text-align: center;
}
li {
  cursor: pointer;
  width: auto;
  height: 100%;
  padding: 5px 15px;
  margin: 0 5px;
  border-top: 1px solid rgba(var(--color-fg-strong-rgb), 0);
  position: relative;
  top: -0.5px;
  display: flex;
  align-items: center;
}
.router-link-exact-active {
  border-top: 1px solid var(--color-primary);
  color: var(--color-primary);
}
</style>
