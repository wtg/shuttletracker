<template>
   <div style="margin-top: 50px;" class="container">
       <are-you-sure @no="shouldDelete = false;" @yes="vehicleDelete" :active="shouldDelete"/>
        <table class="table">
            <thead>
                <tr>
                <th><abbr title="Name">Name</abbr></th>
                <th><abbr title="ID">ID</abbr></th>
                <th><abbr title="Enabled">Enabled</abbr></th>
                <th></th>
                <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="vehicle in vehicles" :key="vehicle.id">
                    <th><router-link :to="'/admin/vehicles/' + String(vehicle.id)">{{vehicle.name}}</router-link></th>
                    <th><p>{{vehicle.id}}</p></th>
                    <th><p>{{vehicle.enabled}}</p></th>
                    <th><button @click="$router.push('/admin/vehicles/' + String(vehicle.id) + '/edit')" class="button">Edit</button></th>
                    <th><button @click="shouldDelete = true; vehicleToDelete = vehicle;" class="button is-danger">Delete</button></th>
                </tr>
                <tr>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th><button @click="$router.push('/admin/vehicles/-1/new')" class="button is-success">New</button></th>
                </tr>
            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import Vehicle from '../../structures/vehicle';
import AdminServiceProvider from '@/structures/serviceproviders/admin.service';
import AreYouSure from '@/components/admin/AreYouSure.vue';

const sp = new AdminServiceProvider();
export default Vue.extend({
    computed: {
        vehicles(): Vehicle[] {
            return this.$store.state.Vehicles;
        },
    },
    data() {
        return{
            shouldDelete: false,
            vehicleToDelete: undefined,
        } as {
            shouldDelete: boolean,
            vehicleToDelete: Vehicle | undefined,
        };
    },
    components: {
        AreYouSure,
    },
    mounted() {
        this.$store.dispatch('grabVehicles');
    },
    methods: {
        vehicleDelete() {
            this.shouldDelete = false;
            if (this.vehicleToDelete !== undefined) {
                AdminServiceProvider.DeleteVehicle((this.vehicleToDelete as Vehicle)).then(() => {
                    this.$store.dispatch('grabVehicles');
                });
            }
        },
    },
});
</script>
