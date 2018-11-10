<template>
    <div style="margin-top: 50px;" class="container">
        <table class="table">
  <thead>
    <tr>
      <th><abbr title="Name">Name</abbr></th>
      <th><abbr title="ID">ID</abbr></th>
      <th><abbr title="Description">Desc.</abbr></th>
      <th><abbr title="Schedule">Sched.</abbr></th>
      <th><abbr title="Enabled">Enabled</abbr></th>
      <th>Edit</th>
    </tr>
  </thead>
  <tbody>
    <tr v-for="route in routes" :key="route.id">
      <th>{{route.name}}</th>
      <td>{{route.id}}</td>
      <td>{{route.description}}</td>
      <td><ul>
            <li v-for="interval in route.schedule" :key="interval.id">{{String(interval)}}</li>
          </ul></td>
      <td>{{route.enabled}}</td>
      <td><button class="button">Edit</button></td>
    </tr>
  </tbody>
</table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import Route from '../../structures/route';

export default Vue.extend({
    name: 'routes',
    computed: {
        routes(): Route[] {
            return this.$store.state.Routes;
        },
    },
    methods: {
        
    },

});
</script>
