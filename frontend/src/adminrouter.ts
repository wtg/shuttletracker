import Vue from 'vue';
import Router from 'vue-router';
import routesAdmin from '@/components/admin/routesAdmin.vue';
import stopsAdmin from '@/components/admin/stopsAdmin.vue';
import routeOverview from '@/components/admin/routeOverview.vue';

Vue.use(Router);

export default new Router({
    mode: 'history',
    routes: [
      {
        path: '/admin/routes/',
        name: 'routes',
        component: routesAdmin,
      },
      {
        path: '/admin/routes/:id',
        name: 'routes',
        component: routeOverview,
      },
      {
        path: '/admin/stops',
        name: 'stops',
        component: stopsAdmin,
      },
    ],
  });
