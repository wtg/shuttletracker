import Vue from 'vue';
import '@/assets/styles.scss';
import Admin from '@/Admin.vue';

Vue.config.productionTip = false;

console.log('here')

/**
 * Declare the main Vue instance with components and vuex store.
 */
new Vue({
    render: (h) => h(Admin),
  }).$mount('#app');
  