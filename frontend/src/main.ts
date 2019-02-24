import Vue from 'vue';
import App from '@/App.vue';
import store from '@/store';
import router from '@/index';
import '@/assets/styles.scss';
import '../node_modules/leaflet/dist/leaflet.css';
import 'typeface-open-sans';

Vue.config.productionTip = false;

/**
 * Declare the main Vue instance with components and vuex store.
 */
new Vue({
  store,
  router,
  render: (h) => h(App),
}).$mount('#app');
