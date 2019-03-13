import Vue from 'vue';
import Router from 'vue-router';
import Public from './components/Public.vue';
import about from './components/about.vue';
import schedules from '@/components/schedules.vue';
import settings from '@/components/settings.vue';
import etas from '@/components/etas.vue';

Vue.use(Router);

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/',
      name: 'map',
      component: Public,
    },
    {
      path: '/about',
      name: 'about',
      component: about,
    },
    {
      path: '/etas',
      name: 'etas',
      component: etas,
    },
    {
      path: '/schedules',
      name: 'schedules',
      component: schedules,
    },
    {
      path: '/settings',
      name: 'settings',
      component: settings,
    },
  ],
});
