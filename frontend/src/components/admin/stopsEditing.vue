<template>
<div class="container">
    <div class="columns">

        <div class="column is-8">
            <div v-if="failure" class="notification is-danger">
                <p>There was an error saving the route.</p>
            </div>
            <div v-if="success" class="notification is-success">
                <p>Route saved successfully</p>
            </div>

            <map-view v-if="false && !creation"></map-view>
            <place-stop v-if="creation"/>
            <!-- <draw-route/> -->
        </div>
        <div class="form-horizontal column" >
            
            <!-- Text input-->
            <div class="field">
            <label class="label" for="Name">Name</label>
            <div class="control">
                <input v-model="stop.name" id="Name" name="Name" type="text" placeholder="Name" class="input" :disabled="!creation">
            </div>
            </div>

            <!-- name, id, description, latitude, longitude, enabled-->

            <!-- Text input-->
            <div class="field">
            <label class="label" for="Description">Description</label>
            <div class="control">
                <input v-model="stop.description" id="Description" name="Description" type="text" placeholder="Description" class="input " :disabled="!creation">
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="Latitude">Latitude</label>
            <div class="control">
                <input v-model="stop.latitude" id="Latitude" name="Latitude" type="text" placeholder="Latitude" class="input " :disabled="!creation">
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="Longitude">Longitude</label>
            <div class="control">
                <input v-model="stop.longitude" id="Longitude" name="Longitude" type="text" placeholder="Longitude" class="input " :disabled="!creation">
            </div>
            </div>



            <!-- Multiple Radios -->
            <!-- <div class="field">
            <label class="label" for="enabled">Show on map</label>
            <div class="control">
                <label class="radio" for="enabled-0">
                <input @click="setEnabled(true);" :checked="route.enabled" type="radio" name="enabled" id="enabled-0" value="enabled">
                yes
                </label>
                <label class="radio" for="enabled-1">
                <input @click="setEnabled(false);" :checked="!route.enabled" type="radio" name="enabled" id="enabled-1" value="disabled">
                no
                </label>
            </div>
            </div>
            <div class="field">
                <label class="label">Schedule Editor</label>
                <schedule-editor v-model="route.schedule" />
            </div> -->
            <div class="field">
            <div class="control">
                <button @click="send" v-if="formValid" :class="{'is-loading': this.sending}" class="button is-info">Save</button>
            </div>
            </div>

        </div>
    </div>
</div>
</template>
<script lang="ts">
import Vue from 'vue';
import mapView from '@/components/admin/map.vue'; // for clicking on a point and inserting a stop
import Stop from '../../structures/stop';
import scheduleEditor from '@/components/admin/scheduleEditor.vue';
import AdminServiceProvider from '../../structures/serviceproviders/admin.service';
import placeStop from '@/components/admin/placeStop.vue';
// import drawRoute from '@/components/admin/drawRoute.vue';
import * as L from 'leaflet';


export default Vue.extend({
    data() {
        return {
            stop: new Stop(-1, '', '', -1, -1, '', ''),
            routePolyLine: undefined,
            sending: false,
            success: false,
            failure: false,
        } as {
            stop: Stop;
            routePolyLine: L.Polyline | undefined;
            sending: boolean;
            success: boolean;
            failure: boolean;
        };
    },
    // },
    components: {
        mapView,
        placeStop,
    },
    // mounted() {
    //     if (this.creation) {
    //         return;
    //     }
    //     // if the routes are not in the state store yet, wait until they are
    //     const el = this;
    //     if (this.$store.getters.getRoutes.length === 0) {
    //         this.$store.subscribe((mutation) => {
    //             if (mutation.type === 'setRoutes') {
    //                 el.grabMyRoute();
    //             }
    //         });
    //     } else {
    //         this.grabMyRoute();
    //     }
    // },
    // watch: {
    //     $route() {
    //         if (this.creation) {
    //             return;
    //         }
    //         this.grabMyRoute();
    //     },
    // },
    computed: {
        formValid(): boolean {
            return this.stop.name !== '';
        },
    },
    methods: {
        send() {
            this.sending = true;
            AdminServiceProvider.NewStop(this.stop).then(() => {
                    this.sending = false;
                    this.success = true;
                    this.$store.dispatch('grabRoutes');
                    setTimeout(() => {
                        this.success = false;
                    }, 2000);
                }).catch(() => {
                    this.failure = true;
                    this.sending = false;
                    setTimeout(() => {
                        this.failure = false;
                    }, 2000);
                });
        },

        // fetch all the fields of the stop
        grabMyStop(){
            for (let i = 0; i < this.$store.getters.getStops.length; i ++) {
                console.log(i);
            }

        }
    },
    props: {
        creation: {
            type: Boolean,
        },
    },
});
</script>
