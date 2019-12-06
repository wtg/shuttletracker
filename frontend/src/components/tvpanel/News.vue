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
    height: 45%;
    width: 100%;
    text-align: center;
    border-top: 3px solid #F8F8F8;
    position:absolute;
  }
  #title {
    margin-top: 20px;
    margin-bottom: 22px;
    font-size: 50px;
    display: inline-block;
    border-bottom: 2px solid #4c4c4c;
  }
  li {
    border: 1px solid rgb(228, 228, 228);
    border-radius: 3px;
    margin-left: 10px;
    margin-right: 10px;
    margin-top: 18px;
    margin-bottom: 18px;
    padding: 10px;
    font-weight: 700;
  }
</style>
