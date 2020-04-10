<template>
    <span>
        <p> Latitude: {{ coordinates.lat }}</p>
        <p> Longitude: {{ coordinates.lng}}</p>
        <div id="map"></div>
    </span>
</template>
<script lang="ts">
import Vue from 'vue';
import * as L from 'leaflet';
import 'leaflet-routing-machine';
import { Stop, StopSVG } from '../../structures/stop';

export default Vue.extend({
    props: {
        stopPoint: {
            type: () => L.marker,
        },
    },
    data() {
        return {
            Map: undefined,
            existingStopMarker: undefined,
            APIKey: 'pk.eyJ1Ijoiamx5b24xIiwiYSI6ImNqNmR4ZTVmejAwaTEzM3FsMmU0d2RmYjIifQ._VUaEMHioVwJIf11PzIqAQ',
            coordinates : new L.LatLng(-1, -1),
        } as {
            Map: L.Map | undefined;
            existingStopMarker: L.Marker | undefined,
            APIKey: string;
            coordinates: L.LatLng | undefined;
        };
    },

    watch: {
        stopPoint() {
            this.renderStops();
        },
    },

    methods: {
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

                    const el = this;
                    this.Map.on('click', (e: any) => {
                            this.coordinates = (e as any).latlng;
                            // sends the coordinates to the form
                            this.$emit('coordinates', this.coordinates);
                    });
                    this.Map.invalidateSize();
                    this.renderStops();
                });
            }
        },

        renderStops(): any {
            // i think this should work, but it doesn't. might need to assign an image.
            if (this.existingStopMarker) {
                if (this.Map !== undefined) {
                    this.existingStopMarker.removeFrom(this.Map);
                }
            }
            if (this.coordinates) {
                this.existingStopMarker = L.marker(this.coordinates, {
                    icon: L.icon({
                        iconUrl: StopSVG,
                        iconSize: [12, 12], // size of the icon
                        iconAnchor: [6, 6], // point of the icon which will correspond to marker's location
                        shadowAnchor: [6, 6], // the same for the shadow
                        popupAnchor: [0, 0], // point from which the popup should open relative to the iconAnchor
                    }),
                });
                if (this.Map !== undefined) {
                    this.existingStopMarker.addTo(this.Map);
                }
            }
        },

    },
    mounted() {
        this.$nextTick(() => {
            this.mountMap();
            this.renderStops();
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
