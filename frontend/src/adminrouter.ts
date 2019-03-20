import Vue from 'vue';
import Router from 'vue-router';
import routesAdmin from '@/components/admin/routesAdmin.vue';
import stopsAdmin from '@/components/admin/stopsAdmin.vue';
import vehiclesAdmin from '@/components/admin/vehiclesAdmin.vue';
import routeOverview from '@/components/admin/routeOverview.vue';
import routeEditing from '@/components/admin/routeEditing.vue';
import stopsEditing from '@/components/admin/stopsEditing.vue';
import vehicleEditing from '@/components/admin/vehicleEditing.vue';
import vehicleOverview from '@/components/admin/vehicleOverview.vue';
import messagesAdmin from '@/components/admin/MessagesAdmin.vue';

Vue.use(Router);

export default new Router({
    mode: 'history',
    routes: [
      {
        path: '/admin',
        redirect: '/admin/routes/',
      },
      {
        path: '/admin/routes/',
        name: 'routes',
        component: routesAdmin,
      },
      {
        path: '/admin/routes/:id',
        name: 'routes view',
        component: routeOverview,
      },
      {
        path: '/admin/routes/:id/edit',
        name: 'routes edit',
        component: routeEditing,
      },
      {
        path: '/admin/routes/:id/new',
        name: 'new route',
        component: routeEditing,
        props: {creation: true},
      },
      {
        path: '/admin/stops/:id/new',
        name: 'new stop',
        component: stopsEditing,
        props: { creation: true },
      },

      {
        path: '/admin/stops',
        name: 'stops',
        component: stopsAdmin,
      },
      {
        path: '/admin/messages',
        name: 'messages',
        component: messagesAdmin,
      },
      {
        path: '/admin/vehicles',
        name: 'vehicles overview',
        component: vehiclesAdmin,
      },
      {
        path: '/admin/vehicles/:id',
        name: 'vehicle overview',
        component: vehicleOverview,
      },
      {
        path: '/admin/vehicles/:id/edit',
        name: 'vehicles',
        component: vehicleEditing,
      },
      {
        path: '/admin/vehicles/:id/new',
        name: 'new vehicle',
        component: vehicleEditing,
        props: {creation: true},
      },
    ],
  });
