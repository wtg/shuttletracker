import Vue from 'vue';
import 'bulma';
import '@/assets/styles.scss';
import Admin from '@/Admin.vue';
import router from '@/adminrouter';
import store from '@/store';
import '../node_modules/leaflet/dist/leaflet.css';
import '../node_modules/leaflet-routing-machine/dist/leaflet-routing-machine.css';

Vue.config.productionTip = false;


/**
 * Declare the main Vue instance with components and vuex store.
 */
new Vue({
    render: (h) => h(Admin),
    router,
    store,
  }).$mount('#app');

