<template>
<div class="container">
    <div class="columns">

        <div class="column is-8">
            <div v-if="failure" class="notification is-danger">
                <p>There was an error saving the stop.</p>
            </div>
            <div v-if="success" class="notification is-success">
                <p>Stop saved successfully</p>
            </div>

            <map-view v-if="false && !creation"></map-view>
            <place-stop v-if="creation" @coordinates="setCoordinates"/>
        </div>
        <div class="form-horizontal column" >
            
            <!-- Text input-->
            <div class="field">
            <label class="label" for="Name">Name</label>
            <div class="control">
                <input v-model="stop.name" id="Name" name="Name" type="text" placeholder="Name" class="input" :disabled="!creation">
            </div>
            </div>

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

            <!-- Submit -->
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
import { Stop } from '../../structures/stop';
import scheduleEditor from '@/components/admin/scheduleEditor.vue';
import AdminServiceProvider from '../../structures/serviceproviders/admin.service';
import placeStop from '@/components/admin/placeStop.vue';
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
    components: {
        mapView,
        placeStop,
    },
    mounted() {
        if (this.creation) {
            return;
        }
        // if the routes are not in the state store yet, wait until they are
        const el = this;
        if (this.$store.getters.getStops.length === 0) {
            this.$store.subscribe((mutation) => {
                if (mutation.type === 'setStops') {
                    el.grabMyStop();
                }
            });
        } else {
            this.grabMyStop();
        }
    },
    watch: {
        $stop() {
            if (this.creation) {
                return;
            }
            this.grabMyStop();
        },
    },
    computed: {
        formValid(): boolean {
            return this.stop.name !== '';
        },
    },
    methods: {
        send() {
            this.sending = true;

            // get most recent data
            // may not be needed
            this.grabMyStop();

            // TODO:
            // Error checking for edit or create
            AdminServiceProvider.NewStop(this.stop).then(() => {
                    this.sending = false;
                    this.success = true;
                    this.$store.dispatch('grabStops');
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

        // method responsible for setting the reactive forms
        setCoordinates(coordinates: L.LatLng) {
            this.stop.latitude = coordinates.lat;
            this.stop.longitude = coordinates.lng;

        },
        // fetch all the fields of the stop
        grabMyStop() {
            for (let i = 0; i < this.$store.getters.getStops.length; i ++) {
                const testStop = this.$store.getters.getStops[i];
                const id = this.stop.id;

                // i have no idea how this statement works
                if (Number(testStop.id) === Number(id)) {
                    console.log(testStop);
                    this.stop.id = testStop.id;
                    this.stop.name = testStop.name;
                    this.stop.description = testStop.description;
                    this.stop.latitude = Number(testStop.latitude);
                    this.stop.longitude = Number(testStop.longitude);
                }
            }
        },
    },
    props: {
        creation: {
            type: Boolean,
        },
    },
});
</script>
