import Vue from 'vue';
import Vuex, { StoreOptions } from 'vuex';
import { StoreState } from './StoreState';
import Route from '@/structures/route';
import InfoServiceProvider from './structures/serviceproviders/info.service';
import Stop from './structures/stop';
import Vehicle from './structures/vehicle';
import * as L from 'leaflet';
import Update from './structures/update';

Vue.use(Vuex);
const InfoService = new InfoServiceProvider();

const store: StoreOptions<StoreState> = {
  state: {
    Routes: [],
    Stops: [],
    Vehicles: [],
  },
  mutations: {
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
          if (Number(vehicle.id) === Number(updates[i].vehicleID)) {

            found = true;
            vehicle.missedUpdates = 0;
            vehicle.setRoute(undefined);
            for (let j = 0; j < state.Routes.length; j ++) {
              if (state.Routes[j].id === updates[i].RouteID) {
                vehicle.setRoute(state.Routes[j]);
                break;
              }
            }
            vehicle.setLatLng(Number(updates[i].lat), Number(updates[i].lng));
            vehicle.showOnMap(true);

            break;
          }
        }
        if (!found) {
          vehicle.missedUpdates ++;
          if (vehicle.missedUpdates > 5) {
            vehicle.showOnMap(false);
          }
        }
      });
    },
  },
  getters: {
    getRoutePolyLines(state: StoreState): L.Polyline[] {
      const arr = new Array<L.Polyline>();
      if (state.Routes !== undefined && state.Routes.length !== 0) {
        state.Routes.forEach((r: Route) => {
          const points = new Array<L.LatLng>();
          if (r.coords !== undefined) {
            r.coords.forEach((p: {lat: number, lng: number}) => {
              points.push(new L.LatLng(p.lat, p.lng));
            });
          }
          const line = new L.Polyline(points, {
            color: r.color,
            weight: r.width,
            opacity: 1,
          });
          arr.push(line);
        });
      }
      return arr;
    },
    getBoundsPolyLine(state: StoreState): L.Polyline {
      const points = new Array<L.LatLng>();
      if (state.Routes !== undefined && state.Routes.length !== 0) {
        state.Routes.forEach((r: Route) => {
          if (r.coords !== undefined) {
            r.coords.forEach((p: {lat: number, lng: number}) => {
              points.push(new L.LatLng(p.lat, p.lng));
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
    grabRotues( {commit} ) {
      InfoService.GrabRoutes().then((ret: Route[]) => commit('setRoutes', ret));
    },
    grabStops( {commit} ) {
      InfoService.GrabStops().then((ret: Stop[]) => commit('setStops', ret));
    },
    grabVehicles( {commit} ) {
      InfoService.GrabVehicles().then((ret: Vehicle[]) => commit('setVehicles', ret));
    },
    grabUpdates( {commit} ) {
      InfoService.GrabUpdates().then((ret: Update[]) => commit('addUpdates', ret));
    },
  },
};

export default new Vuex.Store(store);
