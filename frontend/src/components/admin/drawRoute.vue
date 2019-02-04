<template>
    <span>
        <div>
        <button @click="popWaypoint" class="button is-info undo">Undo</button>
        <p>({{this.RoutingWaypoints.length}}) Waypoints</p>
        </div>
        <div id="map"></div>
    </span>
</template>
<script lang="ts">
import Vue from 'vue';
import * as L from 'leaflet';
import 'leaflet-routing-machine';

export default Vue.extend({
    props: {
        routeLines: {
            type: Array as () => L.Polyline[],
            default: () => [],
        },
    },
    data() {
        return {
            Map: undefined,
            existingRouteLayers: [],
            RoutingControl: undefined,
            APIKey: 'pk.eyJ1Ijoiamx5b24xIiwiYSI6ImNqNmR4ZTVmejAwaTEzM3FsMmU0d2RmYjIifQ._VUaEMHioVwJIf11PzIqAQ',
            drawnRoute: undefined,
            RoutingWaypoints: [],
            RoutePolyLine: undefined,
        } as {
            Map: L.Map | undefined;
            existingRouteLayers: L.Polyline[];
            RoutingControl: L.Routing.Control | undefined;
            APIKey: string;
            drawnRoute: any;
            RoutingWaypoints: any[];
            RoutePolyLine: undefined | L.Polyline;
        };
    },
    watch: {
        routeLines() {
            this.renderRoutes();
        },
    },
    methods: {
        createRoutePoints() {
            const points: any = [];
            const latlngs = (this as any).RoutePolyLine.getLatLngs();
            latlngs.forEach((pt: L.LatLng) => {
                points.push({latitude: pt.lat, longitude: pt.lng});
            });
            this.$emit('points', points);

        },
        mountMap() {
            if (this.Map === undefined) {
                this.$nextTick(() => {
                    this.Map = L.map('map', {
                        zoomControl: false,
                        attributionControl: false,
                    });

                    this.Map.setView([42.728172, -73.678803], 15.3);

                    L.tileLayer('https://stamen-tiles.a.ssl.fastly.net/toner-lite/{z}/{x}/{y}{r}.png', {
                        maxZoom: 17,
                        minZoom: 14,
                    }).addTo(this.Map);
                    this.RoutingControl = L.Routing.control({
                        waypoints: [

                        ],
                        // @ts-ignore
                        createMarker: () => {}, // @ts-ignore
                        // @ts-ignore
                        router: L.Routing.mapbox(this.APIKey),
                        routeWhileDragging: true,
                    });
                    const el = this;

                    this.RoutingControl.on('routeselected', (e) => {
                        (el as any).RoutePolyLine = L.polyline(e.route.coordinates, {color: 'blue'});
                        (el as any).createRoutePoints();
                    });

                    this.RoutingControl.addTo(this.Map);
                    this.RoutingWaypoints = [];
                    this.Map.on('click', (e: any) => {
                            (el as any).RoutingWaypoints.push((e as any).latlng);
                            (el as any).RoutingControl.setWaypoints((el as any).RoutingWaypoints);

                    });
                    this.Map.invalidateSize();
                    this.renderRoutes();

                });
            }
        },
        popWaypoint() {
            (this as any).RoutingWaypoints.splice(-1, 1);
            (this as any).RoutingControl.setWaypoints(this.RoutingWaypoints);
        },
        renderRoutes(): any {
            this.existingRouteLayers.forEach((line) => {
                if (this.Map !== undefined) {
                this.Map.removeLayer(line);
                }
            });
            this.existingRouteLayers = new Array<L.Polyline>();
            const el = this;
            this.routeLines.forEach((line: L.Polyline) => {
                if (el.Map !== undefined) {
                    el.Map.addLayer(line);
                    el.existingRouteLayers.push(line);
                }
            });
        },
    },
    mounted() {
        this.$nextTick(() => {
            this.mountMap();
            this.renderRoutes();
        });
    },
});
</script>

<style lang="scss">
    #map{
        width: 100%;
        height: 500px;
    }
    
    .leaflet-bar {
        display: none !important;
    }

</style>
