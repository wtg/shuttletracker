import Vue from 'vue';
import Vuex, { StoreOptions } from 'vuex';
import { StoreState } from '@/StoreState';
import Route from '@/structures/route';
import InfoServiceProvider from '@/structures/serviceproviders/info.service';
import Stop from '@/structures/stop';
import Vehicle from '@/structures/vehicle';
import * as L from 'leaflet';
import Update from '@/structures/update';
import AdminMessageUpdate from '@/structures/adminMessageUpdate';

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
    adminMessage: undefined,
    online: true,
  },
  mutations: {
    setOnline(state, online: boolean) {
      state.online = online;
    },
    setRoutes(state, routes: Route[]) {
      state.Routes = routes;
    },
    setStops(state, Stops: Stop[]) {
      state.Stops = Stops;
    },
    setVehicles(state, vehicles: Vehicle[]) {
      state.Vehicles = vehicles;
    },
    addUpdates(state, updates: Update[]) {
      const toHide = new Array<Vehicle> ();
      state.Vehicles.forEach((vehicle: Vehicle) => {
        let found = false;
        for (let i = 0; i < updates.length; i++) {
          if (Number(vehicle.id) === Number(updates[i].vehicle_id)) {
            vehicle.lastUpdate = new Date(updates[i].date) ;
            found = true;
            vehicle.speed = Number(updates[i].speed);
            vehicle.setRoute(undefined);
            for (let j = 0; j < state.Routes.length; j ++) {
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
  },
  getters: {
    getRoutePolyLines(state: StoreState): L.Polyline[] {
      const arr = new Array<L.Polyline>();
      if (state.Routes !== undefined && state.Routes.length !== 0) {
        state.Routes.forEach((r: Route) => {
          if (r.enabled) {
            const points = new Array<L.LatLng>();
            if (r.coords !== undefined) {
              r.coords.forEach((p: {latitude: number, longitude: number}) => {
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
          if (r.coords !== undefined) {
            r.coords.forEach((p: {latitude: number, longitude: number}) => {
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
  },
  actions: {
    grabRoutes( {commit} ) {
      InfoService.GrabRoutes().then((ret: Route[]) => commit('setRoutes', ret)).catch(() => {
        commit('setOnline', false);
      });
    },
    grabStops( {commit} ) {
      InfoService.GrabStops().then((ret: Stop[]) => commit('setStops', ret)).catch(() => {
        commit('setOnline', false);
      });
    },
    grabVehicles( {commit} ) {
      InfoService.GrabVehicles().then((ret: Vehicle[]) => commit('setVehicles', ret)).catch(() => {
        commit('setOnline', false);
      });
    },
    grabUpdates( {commit} ) {
      InfoService.GrabUpdates().then((ret: Update[]) => commit('addUpdates', ret)).catch(() => {
        commit('setOnline', false);
      });
    },
    grabAdminMesssage( {commit} ) {
      InfoService.GrabAdminMessage().then((ret: AdminMessageUpdate) => commit('addAdminMessage', ret));
    },
  },

};

export default new Vuex.Store(store);
