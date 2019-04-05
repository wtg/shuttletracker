<template>
  <div id="app">
    <div id="view-wrapper">
      <keep-alive>
        <router-view id="router-view"></router-view>
      </keep-alive>
    </div>  
  </div>
</template>


<script lang="ts">
import Vue from 'vue';
import UserLocationService from '@/structures/userlocation.service';

UserLocationService.getInstance();

export default Vue.extend({
  name: 'app',
});
</script>
<style lang="scss">
#app {
  width: 100%;
  height: 100%;
  overflow: hidden;
}
#view-wrapper {
  height: calc(100% - 40px);
  box-sizing: border-box;
  overflow:hidden;
}
#router-view {
  overflow-y: auto;
  -webkit-overflow-scrolling: touch;
  height: 100%;
  overflow:hidden;
}
</style>
