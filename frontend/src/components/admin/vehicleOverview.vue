<template>
        <div v-if="myVehicle !== null" class="container" style="margin-top: 50px;">
            <div class="level">
                <h1 class="title">
                    {{myVehicle.name}}
                    <button @click="$router.push('/admin/vehicles/' + myVehicle.id + '/edit')" class="button is-info">Edit</button>
                </h1>
            </div>

            <div>
                <p class="has-text-weight-semibold is-size-5">Tracker ID:</p>
                <p>{{myVehicle.tracker_id}}</p>
                <p class="has-text-weight-semibold is-size-5">Enabled:</p>
                <p>{{myVehicle.enabled}}</p>
            </div>

        </div>
</template>
<script lang="ts">
import Vue from 'vue';
import mapView from '@/components/admin/map.vue';
import Vehicle from '@/structures/vehicle';

// This component provides a basic route overview
export default Vue.extend({

    computed: {
        vehicles(): Vehicle[] {
            return this.$store.state.Vehicles;
        },
        myVehicle(): Vehicle | null { // returns the route corresponding to the given id or null if an error occurs
            let ret: Vehicle | null = null;
            this.vehicles.forEach((element: Vehicle) => {
                if (element.id === Number(this.$route.params.id)) {
                    ret = element;
                }
            });
            return ret;
        },

    },
    components: {
        mapView,
    },
});
</script>
<style lang="scss" scoped>

    .routeView{
        display: flex;
        flex-flow: column;
    }

</style>

