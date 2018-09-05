<template>
  <div id="app">
    <Public />
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Public from './components/Public.vue';

export default Vue.extend({
  name: 'app',
  components: {
    Public,
  },
});
</script>

<style lang="scss">

#app{
  width: 100%;
  height: 100%;
}

</style>
