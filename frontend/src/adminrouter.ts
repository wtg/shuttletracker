import Vue from 'vue';
import Router from 'vue-router';
import routesAdmin from '@/components/admin/routesAdmin.vue';
import stopsAdmin from '@/components/admin/stopsAdmin.vue';

Vue.use(Router);

export default new Router({
    mode: 'history',
    routes: [
      {
        path: '/admin/routes',
        name: 'routes',
        component: routesAdmin,
      },
      {
        path: '/admin/stops',
        name: 'stops',
        component: stopsAdmin,
      },
    ],
  });
