<template>
<div class="container">
    <div class="columns">
        <div class="column is-8">
            <map-view v-if="routePolyLine !== undefined" :route-lines="[routePolyLine]"></map-view>
        </div>
        <div class="form-horizontal column" >
            <!-- Text input-->
            <div class="field">
            <label class="label" for="Name">Name</label>
            <div class="control">
                <input v-model="route.name" id="Name" name="Name" type="text" placeholder="Name" class="input" disabled>
                <p v-if="route.name.length === 0" class="help is-danger">Please Enter a name</p>
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="Description">Description</label>
            <div class="control">
                <input v-model="route.description" id="Description" name="Description" type="text" placeholder="Description" class="input " disabled>
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="color">Color</label>
            <div class="control">
                <input v-model="route.color" id="color" name="color" type="text" placeholder="#ff00ff" class="input " disabled>
                <p v-if="!this.colorValid" class="help is-danger">Please set a color</p>
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="width">Width</label>
            <div class="control">
                <input v-model.number="route.width" id="width" name="width" type="number" placeholder="4" class="input " disabled>
                <p v-if="Number(route.width) <= 0" class="help is-danger">Please enter a valid > 0 width</p>
            </div>
            </div>

            <!-- Multiple Radios -->
            <div class="field">
            <label class="label" for="enabled">Enabled/Disabled</label>
            <div class="control">
                <label class="radio" for="enabled-0">
                <input @click="setEnabled(true);" :checked="route.enabled" type="radio" name="enabled" id="enabled-0" value="enabled">
                enabled
                </label>
                <label class="radio" for="enabled-1">
                <input @click="setEnabled(false);" :checked="!route.enabled" type="radio" name="enabled" id="enabled-1" value="disabled">
                disabled
                </label>
            </div>
            </div>
            <div class="field">
                <label class="label">Schedule Editor</label>
                <schedule-editor />
            </div>
            <div class="field">
            <div class="control">
                <button v-if="formValid" class="button is-info">Save</button>
            </div>
            </div>

        </div>
    </div>
</div>
</template>
<script lang="ts">
import Vue from 'vue';
import mapView from '@/components/admin/map.vue';
import Route from '../../structures/route';
import scheduleEditor from '@/components/admin/scheduleEditor.vue';
import * as L from 'leaflet';

export default Vue.extend({
    data(){
        return {
            route: new Route(-1, '', '', true, '#ff00ff', 4, [],[]),
            routePolyLine: undefined,
        } as {
            route: Route;
            routePolyLine: L.Polyline | undefined;
        }
    },
    components: {
        mapView,
        scheduleEditor,
    },
    mounted(){
        // if the routes are not in the state store yet, wait until they are
        let el=this;
        if(this.$store.getters.getRoutes.length === 0){
            this.$store.subscribe((mutation) => {
                if(mutation.type === 'setRoutes'){
                    el.grabMyRoute();
                }
            });
        }else{
            this.grabMyRoute();
        }
    },
    watch: {
        $route(){
            this.grabMyRoute();
        }
    },
    computed: {
        colorValid(): boolean{
            const validColor = new RegExp('^#(?:[0-9a-fA-F]{3}){1,2}$');
            return validColor.test(this.route.color);
        },
        formValid(): boolean {
            return this.route.name !== '' && this.colorValid && this.route.width > 0;
        }
    },
    methods:{
        setEnabled(val: boolean){
            this.route.enabled = val;
        },
        // Grab the route, and set things reactively
        grabMyRoute(){
            for (let i = 0; i < this.$store.getters.getRoutes.length; i ++){
                const testRoute = this.$store.getters.getRoutes[i]
                const id = this.$route.params.id;
                if(Number(testRoute.id) == Number(id)){
                    this.route.id = testRoute.id;
                    this.route.name = testRoute.name;
                    this.route.enabled = testRoute.enabled;
                    this.route.color = testRoute.color;
                    this.route.width = testRoute.width;
                    this.route.description = testRoute.description;
                    this.route.coords = testRoute.coords;
                    this.route.schedule = testRoute.schedule;
                }
                this.routePolyLine = this.$store.getters.getPolyLineByRouteId(testRoute.id);
            }
        }
    }
});
</script>
