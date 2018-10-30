import Vue from 'vue';
import 'bulma';
import '@/assets/styles.scss';
import Admin from '@/Admin.vue';
import router from '@/adminrouter';

Vue.config.productionTip = false;

console.log('here')

/**
 * Declare the main Vue instance with components and vuex store.
 */
new Vue({
    render: (h) => h(Admin),
    router: router,
  }).$mount('#app');
  