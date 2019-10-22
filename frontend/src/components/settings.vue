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

    <table class="table">
        <thead>
            <tr>
            <th><abbr title="Name">Name</abbr></th>
            <th><abbr title="Description">Desc.</abbr></th>
            <th><abbr title="Schedule">Sched.</abbr></th>
            <th><abbr title="Enabled">Enabled</abbr></th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="route in routes" :key="route.id">
            <th><router-link :to='"/admin/routes/" + String(route.id) + "/"' >{{route.name}}</router-link></th>
            <td>{{route.description}}</td>
            <td><ul>
                    <li v-for="interval in route.schedule" :key="interval.id">{{String(interval)}}</li>
                </ul></td>
            <td><b-switch v-model="routeEnabled"></b-switch></td>

            </tr>
        </tbody>

    </table>
    <router-link to="/about">About and privacy policy</router-link>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Route from '../structures/route';

export default Vue.extend({
  name: 'routes',
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
    routes(): Route[] {
        return this.$store.state.Routes;
    },
    routeEnabled: {
        get(): boolean {
            console.log(this.$route.params.name);
            return Boolean(this.$route.params.enabled);
        },
        set(value: boolean) {
            console.log('foo');
            // this.$route.params.enabled = !(this.$route.params.enabled);
            console.log(this.$route.params.enabled);
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
