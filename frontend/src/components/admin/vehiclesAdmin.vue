<template>
   <div style="margin-top: 50px;" class="container">
        <table class="table">
            <thead>
                <tr>
                <th><abbr title="Name">Name</abbr></th>
                <th><abbr title="ID">ID</abbr></th>
                <th><abbr title="Enabled">Enabled</abbr></th>
                <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="vehicle in vehicles" :key="vehicle.id">
                    <th><a>{{vehicle.name}}</a></th>
                    <th><p>{{vehicle.id}}</p></th>
                    <th><p>{{vehicle.enabled}}</p></th>
                    <th><button class="button">Edit</button></th>
                </tr>
            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue'
import Vehicle from '../../structures/vehicle';

export default Vue.extend({
    computed: {
        vehicles(): Vehicle[] {
            return this.$store.state.Vehicle;
        },
    },
})
</script>
