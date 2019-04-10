import Vue from 'vue';
import Router from 'vue-router';
import Public from './components/Public.vue';
import about from './components/about.vue';
import register from '@/components/register.vue';
import schedules from '@/components/schedules.vue';
import settings from '@/components/settings.vue';
import tvpanel from '@/components/tvpanel.vue';
import etas from '@/components/etas.vue';
import Resources from '@/resources';

Vue.use(Router);

export default new Router({
  mode: 'history',
  base: Resources.BasePath,
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
      path: '/register',
      name: 'register',
      component: register,
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
    {
      path: '/tvpanel',
      name: 'tvpanel',
      component: tvpanel,
    },
  ],
});
