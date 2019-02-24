<template>
  <div class="parent">
    <h2 class="title">Settings</h2>
    <hr>
    <div class="field">
      <b-switch>Send position updates</b-switch>
      <p class="help">Use your location to help make Shuttle Tracker more accurate.</p>
    </div>
    <div class="field">
      <b-switch>Bus button</b-switch>
      <p class="help">Place a bus on other users' maps.</p>
    </div>

    <!-- <h3 class="title is-3">About</h3> -->
    <router-link to="/about">About</router-link>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({});
</script>

<style lang="scss" scoped>
.parent {
  padding: 20px;
}
</style>
