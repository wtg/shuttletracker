<template>
  <div class="parent">
    <h1 class="title">Settings</h1>
    <hr />
    <div class="general_section">
      <p class="section_title">General</p>
    </div>

    <table class="table">
      <tbody>
        <tr>
          <th class="route_format">
            <p class="setting_title">Send position updates</p>
            <p
              class="help"
            >Use your location to help make Shuttle Tracker more accurate for everyone. Your location is gathered anonymously while Shuttle Tracker is open.</p>
          </th>
          <td>
            <div class="switches">
              <div class="switch_direction">
                <b-switch v-model="fusionPositionEnabled" v-bind:disabled="geolocationDenied"></b-switch>
              </div>
            </div>
          </td>
        </tr>
        <tr>
          <th class="route_format">
            <p class="setting_title">Bus button</p>
            <p class="help">Place a bus on other users' maps and let others place buses on your map.</p>
          </th>
          <td>
            <div class="switches">
              <div class="switch_direction">
                <b-switch class="internal_switch" v-model="busButtonEnabled"></b-switch>
                <div class="control select is-small is-rounded">
                  <select v-model="busButtonChoice" v-bind:disabled="!busButtonEnabled">
                    <option>🚐</option>
                    <option>🚌</option>
                    <option>🚗</option>
                    <option>🚓</option>
                    <option>🚜</option>
                  </select>
                </div>
              </div>
            </div>
          </td>
        </tr>

        <tr>
          <th class="route_format">
            <p class="setting_title">Estimated times of arrival</p>
            <p
              class="help"
            >Get notifications when a shuttle is likely to arrive at the stop nearest you. Requires access to your location.</p>
            <p class="help">
              <i>Experimental: You’re not allowed to get mad at us if you miss your shuttle.</i>
            </p>
          </th>
          <td>
            <div class="switches">
              <div class="switch_direction">
                <b-switch v-model="etasEnabled"></b-switch>
              </div>
            </div>
          </td>
        </tr>
        <tr>
          <th class="route_format">
            <p class="setting_title">Dark theme</p>
            <p class="help">Turn on dark theme for ShuttleTracker.</p>
          </th>
          <td>
            <div class="switches">
              <div class="switch_direction">
                <b-switch v-model="darkThemeEnabled"></b-switch>
                <div class="control select is-small is-rounded">
                  <select class="is-small" v-model="darkThemeMode" v-bind:disabled="!darkThemeEnabled">
                    <option v-for="mode in darkThemeAllModes" :value="mode.id">{{mode.description}}</option>
                  </select>
                </div>
              </div>
            </div>
          </td>
        </tr>
      </tbody>
    </table>

    <div class="general_section">
      <p class="section_title">Toggle Routes</p>
    </div>
    <table class="table">
      <!-- <thead>
            <tr>
            <th title="Name">Name</th>
            <th title="Enabled">Enabled</th>
            </tr>
      </thead>-->
      <tbody>
        <tr v-for="route in routes" :key="route.id">
          <th class="route_format">{{route.name}}</th>
          <td>
            <b-switch v-bind:value="route.enabled && route.active" v-on:input="routeToggle(route)"></b-switch>
          </td>
        </tr>
      </tbody>
    </table>
    <hr />
    <router-link to="/faq">Frequently Asked Questions</router-link>
    <router-link to="/about">Privacy Policy</router-link>
  </div>
</template>

<script lang="ts">
import Vue from "vue";
import Route from "../structures/route";
import Public from "./Public.vue";
import { DarkTheme, DarkThemeMode } from "../structures/theme";

export default Vue.extend({
  name: "routes",
  computed: {
    busButtonEnabled: {
      get(): boolean {
        return this.$store.state.settings.busButtonEnabled;
      },
      set(value: boolean) {
        this.$store.commit("setSettingsBusButtonEnabled", value);
      }
    },
    busButtonChoice: {
      get(): string {
        return this.$store.state.settings.busButtonChoice;
      },
      set(value: string) {
        this.$store.commit("setSettingsBusButtonChoice", value);
      }
    },
    // darkThemeEnabled is really just a proxy for darkThemeMode. It makes the UI nicer for people who simply want to
    // turn it on or off, while still allowing for more advanced dark theme options via the dropdown.
    darkThemeEnabled: {
      get(): boolean {
        return this.$store.state.settings.darkThemeMode !== "off";
      },
      set(value: boolean) {
        // These assignments will call darkThemeMode.set().
        if (!value) {
          // !enabled -> set mode to off
          this.darkThemeMode = "off";
        } else if (this.darkThemeMode === "off") {
          // enabled -> ensure mode is not off. This prevents overwriting to 'Always' if the user picks a different mode.
          this.darkThemeMode = "on";
        }
      }
    },
    darkThemeMode: {
      get(): string {
        return this.$store.state.settings.darkThemeMode;
      },
      set(value: string) {
        this.$store.commit("setSettingsDarkThemeMode", value);
      }
    },
    darkThemeAllModes(): DarkThemeMode[] {
      return DarkThemeMode.allModes();
    },
    fusionPositionEnabled: {
      get(): boolean {
        return (
          this.$store.state.settings.fusionPositionEnabled &&
          !this.$store.state.geolocationDenied
        );
      },
      set(value: boolean) {
        this.$store.commit("setSettingsFusionPositionEnabled", value);
      }
    },
    geolocationDenied: {
      get(): boolean {
        return this.$store.state.geolocationDenied;
      }
    },
    etasEnabled: {
      get(): boolean {
        return this.$store.state.settings.etasEnabled;
      },
      set(value: boolean) {
        this.$store.commit("setSettingsETAsEnabled", value);
      }
    },
    routes(): Route[] {
      const routes_arr = [];
      for (const stop of this.$store.state.Routes) {
        routes_arr.push(stop);
      }
      return routes_arr;
    }
  },
  methods: {
    routeToggle(thisRoute: Route) {
      // Enable if not currently
      if (!thisRoute.enabled || !thisRoute.active) {
        thisRoute.enabled = true;
        thisRoute.active = true;
      } else {
        // Disable otherwise
        thisRoute.enabled = false;
        thisRoute.active = false;
      }
      this.$store.commit("setRoutes", this.routes);
    }
  }
});
</script>

<style lang="scss" scoped>
.parent {
  padding: 20px;
}

.route_format {
  color: var(--color-fg-strong);
  border-bottom: 1px solid var(--color-divider);
}

td {
  border-bottom: 1px solid var(--color-divider);
}

th {
  border-bottom: 1px solid var(--color-divider);
  font-weight: normal;
}

.switches {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
}
.switch_direction {
  display: flex;
  flex-direction: column;
  align-self: flex-end;
}

.help {
  color: var(--color-fg-normal);
}

.section_title {
  color: var(--color-primary);
  padding: 5px;
}

.setting_title {
  color: var(--color-fg-strong);
}

.options {
  padding: 5px;
}

.field {
  padding: 5px;
  display: inline-block;
}
</style>
