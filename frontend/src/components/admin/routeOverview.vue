<template>
        <div v-if="myRoute !== null" class="container" style="margin-top: 50px;">
            <div class="routeView">
                <h1 class="title">
                    {{myRoute.name}}
                </h1>
                <div class="columns">
                    <div class="column is-8">
                        <map-view :route-lines='[polyLine]' />
                    </div>
                    <div class="column is-4">
                        <div>
                            <p class="has-text-weight-semibold is-size-5">Description:</p>
                            <p>{{myRoute.description}}</p>
                            <p class="has-text-weight-semibold is-size-5">Color:</p>
                            <p>{{myRoute.color}}</p>
                            <p class="has-text-weight-semibold is-size-5">Width:</p>
                            <p>{{myRoute.width}}</p>
                            <p class="has-text-weight-semibold is-size-5">Color:</p>
                            <p>{{myRoute.color}}</p>
                            <p class="has-text-weight-semibold is-size-5">Schedule:</p>
                            <ul>
                                <li v-for="item in myRoute.schedule" :key="item.id">{{String(item)}}</li>
                            </ul>
                        </div>
                    </div>
                </div>
                
            </div>
        </div>
</template>
<script lang="ts">
import Vue from 'vue';
import Route from '../../structures/route';
import mapView from '@/components/admin/map.vue';

// This component provides a basic route overview
export default Vue.extend({
    data() {
        return {
            Map: undefined,
        } as {
            Map: L.Map | undefined;
        };
    },

    computed: {
        routes(): Route[] {
            return this.$store.state.Routes;
        },
        myRoute(): Route | null { // returns the route corresponding to the given id or null if an error occurs
            let ret: Route | null = null;
            this.routes.forEach((element: Route) => {
                if (element.id === Number(this.$route.params.id)) {
                    ret = element;
                }
            });
            return ret;
        },
        polyLine(): L.Polyline | null {
            if (this.myRoute === null) {
                return null;
            } else {
                return this.$store.getters.getPolyLineByRouteId(this.myRoute.id);
            }
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

