<template>
    <div style="margin-top: 50px;" class="container">
        <are-you-sure @no="shouldDelete = false;" @yes="deleteRoute" :active="shouldDelete"/>
        <table class="table">
            <thead>
                <tr>
                <th><abbr title="Name">Name</abbr></th>
                <th><abbr title="ID">ID</abbr></th>
                <th><abbr title="Description">Desc.</abbr></th>
                <th><abbr title="Schedule">Sched.</abbr></th>
                <th><abbr title="Enabled">Enabled</abbr></th>
                <th></th>
                </tr>
            </thead>
            <tbody>
                <tr v-for="route in routes" :key="route.id">
                <th><router-link :to='"/admin/routes/" + String(route.id) + "/"' >{{route.name}}</router-link></th>
                <td>{{route.id}}</td>
                <td>{{route.description}}</td>
                <td><ul>
                        <li v-for="interval in route.schedule" :key="interval.id">{{String(interval)}}</li>
                    </ul></td>
                <td>{{route.enabled}}</td>
                <td><button class="button" @click="$router.push('/admin/routes/' + String(route.id) + '/edit');">Edit</button></td>
                <td><button class="button is-danger" @click="shouldDelete = true; routeToDelete = route;">Delete</button></td>
                
                </tr>
                <tr>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th></th>
                    <th><button @click="$router.push('/admin/routes/-1/new')" class="button is-success">New</button></th>
                </tr>
            </tbody>
        </table>
    </div>
</template>
<script lang="ts">
import Vue from 'vue';
import Route from '../../structures/route';
import AdminServiceProvider from '../../structures/serviceproviders/admin.service';
import AreYouSure from '@/components/admin/AreYouSure.vue';


export default Vue.extend({
    name: 'routes',
    data() {
        return{
            shouldDelete: false,
            routeToDelete: undefined,
        } as {
            shouldDelete: boolean,
            routeToDelete: Route | undefined,
        };
    },
    computed: {
        routes(): Route[] {
            return this.$store.state.Routes;
        },
    },
    methods: {
        deleteRoute() {
            this.shouldDelete = false;
            if (this.routeToDelete !== undefined) {
                AdminServiceProvider.DeleteRoute(this.routeToDelete).then(() => {
                    this.$store.dispatch('grabRoutes');
                });
            }
        },
    },
    components: {
        AreYouSure,
    },

});
</script>
