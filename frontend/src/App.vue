<template>
  <div id="app">
    <keep-alive>
      <router-view id="router-view"></router-view>
    </keep-alive>
    <tab-bar></tab-bar>   
  </div>
</template>


<script lang="ts">
import Vue from 'vue';
import TabBar from '@/components/tabBar.vue';

export default Vue.extend({
  name: 'app',
  components: {
    TabBar,
  },
});
</script>

<style lang="scss">
#app {
  width: 100%;
  height: 100%;
}
#router-view {
  height: calc(100% - 40px);
  overflow: scroll;
  box-sizing: border-box;
}
</style>
