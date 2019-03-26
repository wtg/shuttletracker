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
            <map-view v-if="false && (routePolyLine !== undefined) && !creation" :route-lines="[routePolyLine]"></map-view>
            <draw-route @points="setPoints" v-if="creation"/>
        </div>
        <div class="form-horizontal column" >
            <!-- Text input-->
            <div class="field">
            <label class="label" for="Name">Name</label>
            <div class="control">
                <input v-model="route.name" id="Name" name="Name" type="text" placeholder="Name" class="input" :disabled="!creation">
                <p v-if="route.name.length === 0" class="help is-danger">Please Enter a name</p>
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="Description">Description</label>
            <div class="control">
                <input v-model="route.description" id="Description" name="Description" type="text" placeholder="Description" class="input " :disabled="!creation">
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="color">Color</label>
            <div class="control">
                <input v-model="route.color" id="color" name="color" type="text" placeholder="#ff00ff" class="input " :disabled="!creation">
                <p v-if="!this.colorValid" class="help is-danger">Please set a color</p>
            </div>
            </div>

            <!-- Text input-->
            <div class="field">
            <label class="label" for="width">Width</label>
            <div class="control">
                <input v-model.number="route.width" id="width" name="width" type="number" placeholder="4" class="input " :disabled="!creation">
                <p v-if="Number(route.width) <= 0" class="help is-danger">Please enter a valid > 0 width</p>
            </div>
            </div>

            <!-- Multiple Radios -->
            <div class="field">
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
            </div>
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
import mapView from '@/components/admin/map.vue';
import Route from '../../structures/route';
import scheduleEditor from '@/components/admin/scheduleEditor.vue';
import AdminServiceProvider from '../../structures/serviceproviders/admin.service';
import drawRoute from '@/components/admin/drawRoute.vue';
import * as L from 'leaflet';

// This component is the route editing interface, it creates a new route object on mounted and route change to isolate edits.
// To use route creation mode pass in a :creation="true" prop
export default Vue.extend({
    data() {
        return {
            route: new Route(-1, '', '', true, '#ff00ff', 4, [], [], true, []),
            routePolyLine: undefined,
            sending: false,
            success: false,
            failure: false,
        } as {
            route: Route;
            routePolyLine: L.Polyline | undefined;
            sending: boolean;
            success: boolean;
            failure: boolean;
        };
    },
    components: {
        mapView,
        scheduleEditor,
        drawRoute,
    },
    mounted() {
        if (this.creation) {
            return;
        }
        // if the routes are not in the state store yet, wait until they are
        const el = this;
        if (this.$store.getters.getRoutes.length === 0) {
            this.$store.subscribe((mutation) => {
                if (mutation.type === 'setRoutes') {
                    el.grabMyRoute();
                }
            });
        } else {
            this.grabMyRoute();
        }
    },
    watch: {
        $route() {
            if (this.creation) {
                return;
            }
            this.grabMyRoute();
        },
    },
    computed: {
        colorValid(): boolean {
            const validColor = new RegExp('^#(?:[0-9a-fA-F]{3}){1,2}$');
            return validColor.test(this.route.color);
        },
        formValid(): boolean {
            return this.route.name !== '' && this.colorValid && this.route.width > 0;
        },
    },
    methods: {
        setPoints(points: Array<{latitude: number, longitude: number}>) {
            this.route.points = points;
        },
        send() {
            this.sending = true;
            if (!this.formValid) {
                return;
            }
            if (!this.creation) {
                AdminServiceProvider.EditRoute(this.route).then(() => {
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
            } else {
                AdminServiceProvider.CreateRoute(this.route).then((resp) => {
                    // Handle 500 errors here
                    if (resp.status === 500) {
                        throw new Error('bad response');
                    }
                    this.sending = false;
                    this.success = true;
                    this.$store.dispatch('grabRoutes');
                    setTimeout(() => {
                        this.success = false;
                        this.$router.push('/admin/routes');

                    }, 2000);
                }).catch(() => {
                    this.failure = true;
                    this.sending = false;
                    setTimeout(() => {
                        this.failure = false;
                    }, 2000);
                });
            }

        },
        setEnabled(val: boolean) {
            this.route.enabled = val;
        },
        // Grab the route, and copy it to a new object
        grabMyRoute()    {
            for (let i = 0; i < this.$store.getters.getRoutes.length; i ++) {
                const testRoute = this.$store.getters.getRoutes[i];
                const id = this.$route.params.id;
                if (Number(testRoute.id) === Number(id)) {
                    this.route.id = testRoute.id;
                    this.route.name = testRoute.name;
                    this.route.enabled = testRoute.enabled;
                    this.route.color = testRoute.color;
                    this.route.width = testRoute.width;
                    this.route.description = testRoute.description;
                    this.route.points = testRoute.points;
                    this.route.schedule = testRoute.schedule.slice();
                }
                this.routePolyLine = this.$store.getters.getPolyLineByRouteId(testRoute.id);
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
