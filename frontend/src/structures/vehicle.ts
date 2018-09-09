import * as L from 'leaflet';
import Route from './route';

const ShuttleIcon = require('../assets/shuttle.svg') as string;
const maxMissedUpdatesBeforeHide = 5;

const icon: string = `<?xml version="1.0" encoding="UTF-8"?>
<svg width="52px" height="52px" viewBox="0 0 52 52" version="1.1" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
    <title>shuttle</title>
    <defs></defs>
    <g id="Page-1" stroke="none" stroke-width="1" fill="none" fill-rule="evenodd">
        <g id="shuttle">
            <path d="M51.353,0.914 C51.648,1.218 51.72,1.675 51.532,2.054 L27.532,50.469 C27.362,50.814 27.011,51.025 26.636,51.025 C26.58,51.025 26.524,51.02 26.467,51.01 C26.032,50.936 25.697,50.583 25.643,50.145 L23.098,29.107 L0.835,25.376 C0.402,25.304 0.067,24.958 0.009,24.522 C-0.049,24.086 0.184,23.665 0.583,23.481 L50.218,0.701 C50.603,0.524 51.058,0.609 51.353,0.914 Z" id="Background" fill="#33A7FF"></path>
            <path d="M51.353,0.914 C51.058,0.609 50.603,0.524 50.218,0.701 L0.583,23.481 C0.184,23.665 -0.049,24.086 0.009,24.522 C0.067,24.958 0.402,25.304 0.835,25.376 L23.098,29.107 L25.643,50.145 C25.697,50.583 26.032,50.936 26.467,51.01 C26.524,51.02 26.58,51.025 26.636,51.025 C27.011,51.025 27.362,50.814 27.532,50.469 L51.532,2.054 C51.72,1.675 51.648,1.218 51.353,0.914 Z M27.226,46.582 L24.994,28.125 C24.94,27.685 24.603,27.332 24.166,27.259 L4.374,23.941 L48.485,3.697 L27.226,46.582 Z" id="Shape" fill="#000"></path>
        </g>
    </g>
</svg>`;

/**
 * Vehicle represents a returned vehicle value
 */
export default class Vehicle {
    public id: number;
    public name: string;
    public created: Date;
    public updated: Date;
    public enabled: boolean;
    public marker: L.Marker;
    public lat: number;
    public lng: number;
    public heading: number;
    public speed: number;
    public RouteID: string;
    public shownOnMap: boolean;
    public map: L.Map | undefined;
    public Route: Route | undefined;
    public missedUpdates: number;

    constructor(id: number, name: string, created: Date, updated: Date, enabled: boolean) {
        this.id = id;
        this.name = name;
        this.created = created;
        this.updated = updated;
        this.enabled = enabled;
        this.lat = 0;
        this.lng = 0;
        this.heading = 0;
        this.speed = 0;
        this.RouteID = '';
        this.marker = new L.Marker([this.lat, this.lng], {
            icon: L.icon({
                iconUrl: 'data:image/svg+xml;base64,' + btoa(icon),
                iconSize:     [32, 32], // size of the icon
                iconAnchor:   [16, 16], // point of the icon which will correspond to marker's location
                popupAnchor:  [0, 0],   // point from which the popup should open relative to the iconAnchor
            }),
            zIndexOffset: 1000,
        });
        this.shownOnMap = false;
        this.map = undefined;
        this.Route = undefined;
        this.missedUpdates = 0;
    }

    public addToMap(map: L.Map) {
        if (this.map === undefined) {
            this.map = map;
        }
        this.marker.addTo(map);
        this.shownOnMap = true;
    }

    public showOnMap(show: boolean) {
        if (show) {
            if (!this.shownOnMap && this.map !== undefined) {
                this.addToMap(this.map);
                this.shownOnMap = true;
            }
        } else {
            if (this.shownOnMap && this.map !== undefined) {
                this.marker.removeFrom(this.map);
                this.shownOnMap = false;
            }
        }
    }

    public setRoute(r: Route | undefined) {
        if (r === undefined) {
            this.marker.setIcon(L.icon({
                iconUrl: 'data:image/svg+xml;base64,' + btoa(icon.replace('#33A7FF', '#FFFFFF')),
                iconSize:     [32, 32], // size of the icon
                iconAnchor:   [16, 16], // point of the icon which will correspond to marker's location
                popupAnchor:  [0, 0],   // point from which the popup should open relative to the iconAnchor
            }));
            return;
        }
        this.Route = r;
        this.RouteID = r.id;
        this.marker.setIcon(L.icon({
            iconUrl: 'data:image/svg+xml;base64,' + btoa(icon.replace('#33A7FF', '#FFF')),
            iconSize:     [32, 32], // size of the icon
            iconAnchor:   [16, 16], // point of the icon which will correspond to marker's location
            popupAnchor:  [0, 0],   // point from which the popup should open relative to the iconAnchor
        }));
    }

    public setLatLng(lat: number, lng: number) {
        this.lat = lat;
        this.lng = lng;
        this.marker.setLatLng([this.lat, this.lng]);
    }

    public removeFromMap(map: L.Map) {
        map.removeLayer(this.marker);
    }
}
