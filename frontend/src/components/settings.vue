<template>
  <div class="parent">
    <h2 class="title">Settings</h2>
    <hr>
    <div class="field">
      <b-switch v-model="fusionPositionEnabled">Send position updates</b-switch>
      <p class="help">Use your location to help make Shuttle Tracker more accurate. When enabled your location is collected anonymously while Shuttle Tracker is open.</p>
    </div>
    <div class="field">
      <b-switch v-model="busButtonEnabled">Bus button</b-switch>
      <p class="help">Place a bus on other users' maps.</p>
    </div>

    <!-- <h3 class="title is-3">About</h3> -->
    <router-link to="/about">About and privacy policy</router-link>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';

export default Vue.extend({
  computed: {
    busButtonEnabled: {
      get(): boolean {
        return this.$store.state.settings.busButtonEnabled;
      },
      set(value: boolean) {
        this.$store.commit('setSettingsBusButtonEnabled', value);
      },
    },
    fusionPositionEnabled: {
      get(): boolean {
        return this.$store.state.settings.fusionPositionEnabled;
      },
      set(value: boolean) {
        this.$store.commit('setSettingsFusionPositionEnabled', value);
      },
    },
  },
});
</script>

<style lang="scss" scoped>
.parent {
  padding: 20px;
}
</style>
