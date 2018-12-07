<template>
  <div id="admin">
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a class="navbar-item" href="https://shuttles.rpi.edu">
          <p>Shuttle Tracc - <span class="has-text-danger">Fleet Management</span></p>
        </a>
      </div>

      <div id="navbarBasicExample" class="navbar-menu">
        <div class="navbar-start">
          <router-link to="/admin/routes" class="navbar-item">
            Routes
          </router-link>
          <router-link to="/admin/stops" class="navbar-item">
            Stops
          </router-link>
          <router-link to="/admin/vehicles" class="navbar-item">
            Vehicles
          </router-link >
          <router-link to="/admin/messages" class="navbar-item">
            Messages
          </router-link>
        </div>         

        <div class="navbar-end">
          <div class="navbar-item">
            <div class="buttons">
              <a href="/admin/logout" class="button is-light">
                Log out
              </a>
            </div>
          </div>
        </div>
      </div>
    </nav>
    <router-view>
    </router-view>
  </div>
</template>


<script lang="ts">
import Vue from 'vue';
import routesAdmin from '@/components/admin/routesAdmin.vue';

export default Vue.extend({
  name: 'admin',
  components: {
    routesAdmin,
  },
  mounted() {
    this.$store.dispatch('grabRoutes');
    this.$store.dispatch('grabVehicles');
    this.$store.dispatch('grabStops');
  },
});
</script>

<style lang="scss">

#app{
  width: 100%;
  height: 100%;
}

</style>
<
