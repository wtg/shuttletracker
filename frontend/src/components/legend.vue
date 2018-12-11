<template>
    <div class="legend">
        <p>Info Legend</p>
    </div>
</template>
        
<script lang="ts">
// This component is the legend that is displayed on the bottom left corner

import Vue from 'vue';
export default Vue.extend({

});
</script>
<style lang="scss">
    .legend{
      box-shadow: rgba(0, 0, 0, 0.8) 0px 1px 1px;
      border-radius: 5px;
      background-color: rgba(255, 255, 255, .9);
      padding: 5px;
      bottom: 25px;
} 
</style>