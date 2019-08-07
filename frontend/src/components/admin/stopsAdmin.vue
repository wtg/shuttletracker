<template>
    <div style="margin-top: 50px;" class="container">
        <are-you-sure @no="shouldDelete = false;" @yes="deleteStop" :active="shouldDelete"/>
        <table class="table">
            <thead>
                <tr>
                    <th><abbr title="Name">Name</abbr></th>
                    <th><abbr title="ID">ID</abbr></th>
                    <th><abbr title="Desc">Description</abbr></th>
                    <th><abbr title="Lat.">Latitude</abbr></th>
                    <th><abbr title="Lng.">Longitude</abbr></th>
                    <th><abbr title="Enabled">Enabled</abbr></th>
                    <th></th>
                    <th></th>
                </tr>
                <tr v-for="stop in stops" :key="stop.id">
                    <th>{{stop.name}}</th>
                    <th>{{stop.id}}</th>
                    <th>{{stop.description}}</th>
                    <th>{{stop.latitude}}</th>
                    <th>{{stop.longitude}}</th>
                    <th></th>
                    <th><button @click="$router.push('/admin/stops/' + String(stop.id) + '/edit')" class="button">Edit</button></th>
                    <td><button class="button is-danger" @click="shouldDelete = true; stopToDelete = stop;">Delete</button></td>
                </tr>
            </thead>
            <tbody>

                <tr>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th><button @click="$router.push('/admin/stops/-1/new')" class="button is-success">New</button></th>
                </tr>

            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import { Stop } from '@/structures/stop';
import AreYouSure from '@/components/admin/AreYouSure.vue';
import AdminServiceProvider from '../../structures/serviceproviders/admin.service';
export default Vue.extend({
    name: 'stops',
    data() {
        return{
            shouldDelete: false,
            stopToDelete: undefined,
        } as {
            shouldDelete: boolean,
            stopToDelete: Stop | undefined,
        };
    },
    computed: {
        stops(): Stop[] {
            return this.$store.state.Stops;
        },
    },
    methods: {
        deleteStop() {
            this.shouldDelete = false;
            if (this.stopToDelete !== undefined) {
                AdminServiceProvider.DeleteStop(this.stopToDelete).then(() => {
                    this.$store.dispatch('grabStops');
                });
            }
        },
    },
    components: {
        AreYouSure,
    },

});
</script>
