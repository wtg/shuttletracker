<template>
  <div id="news">
    <div id="title">News</div>
  </div>
</template>

<script lang="ts">

// This component is the News component for the tv panel
import Vue from 'vue';

export default Vue.extend({
  name: 'News',
});
</script>

<style lang="scss" scoped>

#news {
  height: 50%;
  width: 100%;
  text-align: center;
  border-top: 1.5px solid #F8F8F8;
}

#title {
  margin-top: 20px;
  font-size: 50px;
  display: inline-block;
  border-bottom: 2px solid #4c4c4c;
}

</style>
