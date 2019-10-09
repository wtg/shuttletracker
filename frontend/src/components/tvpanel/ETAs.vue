<template>
  <div id="etas">
    <div id="title">ETAs</div>
  </div>
</template>

<script lang="ts">

// This component is the ETAs component for the tv panel
import Vue from 'vue';

export default Vue.extend({
  name: 'ETAs',
});
</script>

<style lang="scss" scoped>
  #etas {
    height: 50%;
    width: 100%;
    text-align: center;
    border-bottom: 1.5px solid #F8F8F8;
  }

  #title {
    margin-top: 20px;
    font-size: 50px;
    display: inline-block;
    border-bottom: 2px solid #4c4c4c;
  }

</style>
