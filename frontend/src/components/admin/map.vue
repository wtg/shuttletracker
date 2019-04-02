<template>
    <div id="map"></div>
</template>
<script lang="ts">
import Vue from 'vue';
import * as L from 'leaflet';

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

        } as {
            Map: L.Map | undefined;
            existingRouteLayers: L.Polyline[];

        };
    },
    watch: {
        routeLines() {
            this.renderRoutes();
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

                    this.Map.addControl(L.control.attribution({
                        position: 'bottomright',
                        prefix: '',
                    }));
                    L.tileLayer('https://stamen-tiles.a.ssl.fastly.net/toner-lite/{z}/{x}/{y}{r}.png', {
                        attribution: 'Map tiles by <a href="http://stamen.com">Stamen Design</a>, ' +
                                    'under <a href="http://creativecommons.org/licenses/by/3.0">CC BY 3.0</a>. ' +
                                    'Data by <a href="http://openstreetmap.org">OpenStreetMap</a>, under ' +
                                    '<a href="http://www.openstreetmap.org/copyright">ODbL</a>.',
                        maxZoom: 17,
                        minZoom: 14,
                    }).addTo(this.Map);

                    this.Map.invalidateSize();
                    this.renderRoutes();

                });
            }
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

<style lang="scss" scoped>
    #map{
        width: 100%;
        height: 500px;
    }
</style>
