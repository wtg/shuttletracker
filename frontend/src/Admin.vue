<template>
  <div id="admin">
    <nav class="navbar" role="navigation" aria-label="main navigation">
      <div class="navbar-brand">
        <a class="navbar-item" href="https://shuttles.rpi.edu">
          <p>Shuttle iTrak - <span class="has-text-danger">Fleet Management</span></p>
        </a>
      </div>

      <div id="navbarBasicExample" class="navbar-menu">
        <div class="navbar-start">
          <router-link to="radmin" class="navbar-item">
            Routes
          </router-link>
          <router-link to="sadmin" class="navbar-item">
            Stops
          </router-link>
          <a class="navbar-item">
            Vehicles
          </a>
          <a class="navbar-item">
            Messages
          </a>
        </div>         

        <div class="navbar-end">
          <div class="navbar-item">
            <div class="buttons">
              <a class="button is-light">
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
});
</script>

<style lang="scss">

#app{
  width: 100%;
  height: 100%;
}

</style>
<
