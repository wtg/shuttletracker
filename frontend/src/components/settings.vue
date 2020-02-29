<template>
  <div class="parent">
    <h1 class="title">Settings</h1>
    <hr>
    <div class="field">
      <b-switch v-model="fusionPositionEnabled" v-bind:disabled="geolocationDenied">Send position updates</b-switch>
      <p class="help">Use your location to help make Shuttle Tracker more accurate for everyone. Your location is gathered anonymously while Shuttle Tracker is open.</p>
    </div>
    <div class="field">
    <b-switch v-model="busButtonEnabled">  Bus button <span>
      <div class="field">
          <div class="control">
              <div class="select is-small">
                  <div class="select is-danger">
                      <div class="select is-rounded">
                          <select v-model="busButtonChoice">
                              <option>🚐</option>
                              <option>🚌</option>
                              <option>🚗</option>
                              <option>🚓</option>
                              <option>🚜</option>
                          </select>
                      </div>
                  </div>
              </div>
          </div>
      </div>
      </span>
      </b-switch>
      <p class="help">Place a bus on other users' maps and let others place buses on your map.</p>
    </div>

    <b-field v-bind:message="['Get notifications when a shuttle is likely to arrive at the stop nearest you. Requires access to your location.', '<i>Warning: this feature is experimental. You’re not allowed to get mad at us if you miss your shuttle.</i>']">
      <b-switch v-model="etasEnabled">Estimated times of arrival</b-switch>
    </b-field>

    <div class="field">
      <b-switch v-model="darkThemeEnabled">Dark theme
        <div class="control select is-small is-danger is-rounded">
          <select v-model="darkThemeMode" v-bind:disabled="!darkThemeEnabled">
            <option v-for="mode in darkThemeAllModes" :value="mode.id">
              {{mode.description}}
            </option>
          </select>
        </div>
      </b-switch>
      <p class="help">Change the theme to Dark across all pages.</p>
    </div>

    <table class="table">
        <thead>
            <tr>
            <th><abbr title="Name">Name</abbr></th>
            <th><abbr title="Enabled">Enabled</abbr></th>
            </tr>
        </thead>
        <tbody>
            <tr v-for="route in routes" :key="route.id">
            <th class = "route_format">{{route.name}}</th>
            <td><b-switch v-bind:value="route.enabled && route.active" v-on:input="routeToggle(route)"></b-switch></td>
            </tr>
        </tbody>
    </table>
    <router-link to="/faq">Frequently Asked Questions</router-link>
    <router-link to="/about">Privacy Policy</router-link>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import Route from '../structures/route';
import Public from './Public.vue';
import {DarkTheme, DarkThemeMode} from '@/structures/theme';

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
    busButtonChoice: {
      get(): string {
        return this.$store.state.settings.busButtonChoice;
      },
      set(value: string) {
        this.$store.commit('setSettingsBusButtonChoice', value);
      },
    },
    // darkThemeEnabled is really just a proxy for darkThemeMode. It makes the UI nicer for people who simply want to
    // turn it on or off, while still allowing for more advanced dark theme options via the dropdown.
    darkThemeEnabled: {
      get(): boolean {
        return this.$store.state.settings.darkThemeMode !== 'off';
      },
      set(value: boolean) {
        // These assignments will call darkThemeMode.set().
        if (!value) {
          // !enabled -> set mode to off
          this.darkThemeMode = 'off';
        } else if (this.darkThemeMode === 'off') {
          // enabled -> ensure mode is not off. This prevents overwriting to 'Always' if the user picks a different mode.
          this.darkThemeMode = 'on';
        }
      },
    },
    darkThemeMode: {
      get(): string {
        return this.$store.state.settings.darkThemeMode;
      },
      set(value: string) {
        this.$store.commit('setSettingsDarkThemeMode', value);
      },
    },
    darkThemeAllModes(): DarkThemeMode[] {
      return DarkThemeMode.allModes();
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
        const routes_arr = [];
        for (const stop of this.$store.state.Routes) {
          routes_arr.push(stop);
        }
        return routes_arr;
    },
  },
  methods: {
      routeToggle(thisRoute: Route) {
        // Enable if not currently
        if (!thisRoute.enabled || !thisRoute.active) {
            thisRoute.enabled = true;
            thisRoute.active = true;
        } else {  // Disable otherwise
            thisRoute.enabled = false;
            thisRoute.active = false;
        }
        this.$store.commit('setRoutes', this.routes);
      },
  },
});
</script>

<style lang="scss" scoped>
.parent {
  padding: 20px;
}

.route_format {
  color: red;
  font-weight: bold;
}
.field {
    display: inline-block;
}
</style>
