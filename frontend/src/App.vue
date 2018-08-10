<template>
  <div id="app">
    <titlebar />
  </div>
</template>

<script lang="ts">
import { Component, Vue } from 'vue-property-decorator';
import titlebar from './components/titlebar.vue';
@Component({
  components: {
    titlebar,
  },
})

export default class App extends Vue {}
</script>

<style lang="scss">

</style>
