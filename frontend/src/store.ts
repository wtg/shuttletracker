import Vue from 'vue';
import Vuex, { StoreOptions } from 'vuex';
import { StoreState } from '@/StoreState';
import Route from '@/structures/route';
import InfoServiceProvider from '@/structures/serviceproviders/info.service';
import { Stop } from '@/structures/stop';
import Vehicle from '@/structures/vehicle';
import ETA from '@/structures/eta';
import * as L from 'leaflet';
import Update from '@/structures/update';
import AdminMessageUpdate from '@/structures/adminMessageUpdate';
import { stat } from 'fs';

Vue.use(Vuex);
const InfoService = new InfoServiceProvider();

/**
 * The vuex store will contain all state data for the application, mutations allow modification of the store data
 * getters allow reactive getting with modifications, and actions allow for asyncronous data collection. Only actions
 * should be used to make API calls in order to force things to be done in a clean manner.
 */

const store: StoreOptions<StoreState> = {
  state: {
    Routes: [],
    Stops: [],
    Vehicles: [],
    etas: [],
    adminMessage: undefined,
    online: true,
    settings: {
      busButtonEnabled: false,
      etasEnabled: false,
      fusionPositionEnabled: true,
    },
    geolocationDenied: false,
  },
  mutations: {
    setOnline(state, online: boolean) {
      state.online = online;
    },
    setRoutes(state, routes: Route[]) {
      state.Routes = routes;

      // also ensure that vehicles consider being on any newly-returned routes
      state.Vehicles.forEach((vehicle: Vehicle) => {
        state.Routes.forEach((route: Route) => {
          if (vehicle.RouteID === route.id) {
            vehicle.setRoute(route);
            return;
          }
        });
      });
    },
    setStops(state, Stops: Stop[]) {
      state.Stops = Stops;
    },
    setVehicles(state, vehicles: Vehicle[]) {
      state.Vehicles = vehicles;

      // also ensure that vehicles consider being on any known routes
      state.Vehicles.forEach((vehicle: Vehicle) => {
        state.Routes.forEach((route: Route) => {
          if (vehicle.RouteID === route.id) {
            vehicle.setRoute(route);
            return;
          }
        });
      });
    },
    addUpdates(state, updates: Update[]) {
      const toHide = new Array<Vehicle>();
      state.Vehicles.forEach((vehicle: Vehicle) => {
        let found = false;
        for (let i = 0; i < updates.length; i++) {
          if (Number(vehicle.id) === Number(updates[i].vehicle_id)) {
            vehicle.lastUpdate = new Date(updates[i].time);
            found = true;
            vehicle.speed = Number(updates[i].speed);
            vehicle.setRoute(undefined);
            for (let j = 0; j < state.Routes.length; j++) {
              if (state.Routes[j].id === updates[i].route_id) {
                vehicle.setRoute(state.Routes[j]);
                break;
              }
            }
            vehicle.setLatLng(Number(updates[i].latitude), Number(updates[i].longitude));
            vehicle.setHeading(Number(updates[i].heading));
            vehicle.showOnMap(true);

            break;
          }
        }
        if (!found) {
          vehicle.showOnMap(false);
        }
      });
    },
    addAdminMessage(state, message: AdminMessageUpdate) {
      state.adminMessage = message;
    },
    updateETAs(state, { vehicleID, etas }) {
      // remove this vehicle's ETAs and any ETA that has expired
      const now = new Date();
      for (let i = state.etas.length - 1; i >= 0; i--) {
        const eta = state.etas[i];
        if (eta.vehicleID === vehicleID || eta.eta < now) {
          state.etas.splice(i, 1);
        }
      }

      // store new ETAs
      for (const eta of etas) {
        state.etas.push(eta);
      }
    },
    initializeSettings(state) {
      const savedSettings = localStorage.getItem('st_settings');
      if (!savedSettings) {
        return;
      }
      state.settings = Object.assign(state.settings, JSON.parse(savedSettings));
    },
    setSettingsBusButtonEnabled(state, value: boolean) {
      state.settings.busButtonEnabled = value;
      localStorage.setItem('st_settings', JSON.stringify(state.settings));
    },
    setSettingsETAsEnabled(state, value: boolean) {
      state.settings.etasEnabled = value;
      localStorage.setItem('st_settings', JSON.stringify(state.settings));
    },
    setSettingsFusionPositionEnabled(state, value: boolean) {
      state.settings.fusionPositionEnabled = value;
      localStorage.setItem('st_settings', JSON.stringify(state.settings));
    },
    setGeolocationDenied(state, value: boolean) {
      state.geolocationDenied = value;
    },
    setRoutesOnStops(state) {
      // set any routes on existing stops
      state.Stops.forEach((stop: Stop) => {
        state.Routes.forEach((route: Route) => {
          if (route.containsStop(stop.id) && route.active) {
            stop.addRoute(route);
          }
        });
      });
    },
  },
  getters: {
    getBusButtonVisible(state: StoreState, getters): boolean {
      return getters.getBusButtonShowBuses && !state.geolocationDenied;
    },
    getBusButtonShowBuses(state: StoreState): boolean {
      if (state.adminMessage === undefined) {
        return state.settings.busButtonEnabled;
      }
      return state.settings.busButtonEnabled && !state.adminMessage.enabled;
    },
    getPolyLineByRouteId: (state) => (id: number): L.Polyline | undefined => {
      const arr = new Array<L.Polyline>();
      let ret;

      if (state.Routes !== undefined && state.Routes.length !== 0) {
        state.Routes.forEach((r: Route) => {
          if (r.enabled) {
            const points = new Array<L.LatLng>();
            if (r.points !== undefined) {
              r.points.forEach((p: { latitude: number, longitude: number }) => {
                points.push(new L.LatLng(p.latitude, p.longitude));
              });
            }
            const line = new L.Polyline(points, {
              color: r.color,
              weight: r.width,
              opacity: 1,
            });
            if (r.id === id) {
              ret = line;
            }
          }
        });
      }
      return ret;
    },
    getRoutePolyLines(state: StoreState): L.Polyline[] {
      const arr = new Array<L.Polyline>();
      if (state.Routes !== undefined && state.Routes.length !== 0) {
        state.Routes.forEach((r: Route) => {
          if (r.shouldShow()) {
            const points = new Array<L.LatLng>();
            if (r.points !== undefined) {
              r.points.forEach((p: { latitude: number, longitude: number }) => {
                points.push(new L.LatLng(p.latitude, p.longitude));
              });
            }
            const line = new L.Polyline(points, {
              color: r.color,
              weight: r.width,
              opacity: 1,
            });
            arr.push(line);
          }
        });
      }
      return arr;
    },
    getBoundsPolyLine(state: StoreState): L.Polyline {
      const points = new Array<L.LatLng>();
      if (state.Routes !== undefined && state.Routes.length !== 0) {
        state.Routes.forEach((r: Route) => {
          if (r.shouldShow() && r.points !== undefined) {
            r.points.forEach((p: { latitude: number, longitude: number }) => {
              points.push(new L.LatLng(p.latitude, p.longitude));
            });
          }
        });

      }
      const line = new L.Polyline(points, {
        opacity: 1,
      });
      return line;
    },
    getRoutes(state: StoreState): Route[] {
      return state.Routes;
    },
    getStops(state: StoreState): Stop[] {
      return state.Stops;
    },
    getVehicles(state: StoreState): Vehicle[] {
      return state.Vehicles;
    },
  },
  actions: {
    grabRoutes({ commit }) {
      InfoService.GrabRoutes().then((ret: Route[]) => commit('setRoutes', ret)).catch(() => {
        commit('setOnline', false);
      });
    },
    grabStops({ commit }) {
      InfoService.GrabStops().then((ret: Stop[]) => commit('setStops', ret)).catch(() => {
        commit('setOnline', false);
      });
    },
    grabVehicles({ commit }) {
      InfoService.GrabVehicles().then((ret: Vehicle[]) => commit('setVehicles', ret)).catch(() => {
        commit('setOnline', false);
      });
    },
    grabUpdates({ commit }) {
      InfoService.GrabUpdates().then((ret: Update[]) => commit('addUpdates', ret)).catch(() => {
        commit('setOnline', false);
      });
    },
    grabAdminMesssage({ commit }) {
      InfoService.GrabAdminMessage().then((ret: AdminMessageUpdate) => commit('addAdminMessage', ret));
    },
  },

};

export default new Vuex.Store(store);
