<template>
  <div id="app">
    <keep-alive>
      <router-view />

    </keep-alive>
  </div>
</template>


<script lang="ts">
import Vue from 'vue';
import Public from './components/Public.vue';

export default Vue.extend({
  name: 'app',
});
</script>

<style lang="scss">

#app{
  width: 100%;
  height: 100%;
}

</style>
