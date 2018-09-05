import Vue from 'vue';
import App from './App.vue';
import store from './store';
import './assets/styles.scss';
import '../node_modules/leaflet/dist/leaflet.css';

Vue.config.productionTip = false;

new Vue({
  store,
  render: (h) => h(App),
}).$mount('#app');
