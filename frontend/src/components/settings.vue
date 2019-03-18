<template>
  <div class="parent">
    <h1 class="title">Settings</h1>
    <hr>
    <div class="field">
      <b-switch v-model="fusionPositionEnabled" v-bind:disabled="geolocationDenied">Send position updates</b-switch>
      <p class="help">Use your location to help make Shuttle Tracker more accurate for everyone. Your location is gathered anonymously while Shuttle Tracker is open.</p>
    </div>
    <div class="field">
      <b-switch v-model="busButtonEnabled">Bus button</b-switch>
      <p class="help">Place a bus on other users' maps and let others place buses on your map.</p>
    </div>

    <b-field v-bind:message="['Get notifications when a shuttle is likely to arrive at the stop nearest you. Requires access to your location.', '<i>Warning: this feature is experimental. You’re not allowed to get mad at us if you miss your shuttle.</i>']">
      <b-switch v-model="etasEnabled">Estimated times of arrival</b-switch>
    </b-field>

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
        return this.$store.state.settings.fusionPositionEnabled && !this.$store.state.geolocationDenied;
      },
      set(value: boolean) {
        this.$store.commit('setSettingsFusionPositionEnabled', value);
      },
    },
    geolocationDenied: {
      get(): boolean {
        return this.$store.state.geolocationDenied;
      },
    },
    etasEnabled: {
      get(): boolean {
        return this.$store.state.settings.etasEnabled;
      },
      set(value: boolean) {
        this.$store.commit('setSettingsETAsEnabled', value);
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
