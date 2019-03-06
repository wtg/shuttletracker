<<<<<<< HEAD
import Vue from 'vue';
import Router from 'vue-router';
import Public from './components/Public.vue';
import about from './components/about.vue';
import tvcomponent from './components/tvpanel.vue';

Vue.use(Router);

export default new Router({
    mode: 'history',
    routes: [
      {
        path: '/',
        name: 'main',
        component: Public,
      },
      {
        path: '/about',
        name: 'about',
        component: about,
      },
      {
        path: '/tvpanel',
        name: 'TV Component',
        component: tvcomponent,
      }
    ],
  });
=======
import Vue from 'vue';
import Router from 'vue-router';
import Public from './components/Public.vue';
import about from './components/about.vue';
import schedules from '@/components/schedules.vue';
import settings from '@/components/settings.vue';

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
>>>>>>> 727997b041c28c0ba29e04b832e8cd36605f18af
