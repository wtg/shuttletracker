import Vue from 'vue';
import Buefy from 'buefy';

// @ts-ignore
import VueAnalytics from 'vue-analytics';

// styles
import '@fortawesome/fontawesome-free/css/fontawesome.css';
import '@fortawesome/fontawesome-free/css/solid.css'; // only include the specific fontawesome icons that we use
import 'typeface-open-sans';
import 'leaflet/dist/leaflet.css';
import 'buefy/dist/buefy.css';

import '@/assets/vars.scss';
import '@/assets/styles.scss';

import App from '@/App.vue';
import store from '@/store';
import router from '@/index';

Vue.use(Buefy);
Vue.config.productionTip = false;

Vue.use(VueAnalytics, {
  id: 'UA-28203673-6',
  autoTracking: {
    exception: true,
  },
  router,
});

/**
 * Declare the main Vue instance with components and vuex store.
 */
new Vue({
  store,
  router,
  render: (h) => h(App),
  beforeCreate() {
    this.$store.commit('initializeSettings');
  },
}).$mount('#app');

