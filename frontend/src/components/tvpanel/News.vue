<template>
  <div id="news">
    <div id="title">News</div>
    <ul>
      <li>Rensselaer community unsettled by police presence on campus</li>
      <li>Period products absent from campus, Student Government takes action</li>
      <li>Indian dance clubs unite</li>
    </ul>
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
  padding: 10px;
  border-top: 1.5px solid #F8F8F8;
}

#title {
  font-size: 50px;
  display: inline-block;
  margin-bottom: 10px;
  border-bottom: 2px solid #4c4c4c;
}

li {
  border: 1px solid rgb(228, 228, 228);
  border-radius: 3px;
  margin: 10px;
  padding: 10px;
}

</style>
