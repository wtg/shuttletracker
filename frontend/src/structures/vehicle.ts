import * as L from 'leaflet';
import Route from '@/structures/route';
import Location from '@/structures/location';
import getMarkerString from '@/structures/leaflet/rotatedMarker';
import getCardinalDirection from '@/structures/cardinalDirection';
import 'leaflet-rotatedmarker';

const ShuttleIcon = require('@/assets/shuttle.svg') as string;
const maxMissedUpdatesBeforeHide = 5;


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
    public RouteID: number | null;
    public shownOnMap: boolean;
    public map: L.Map | undefined;
    public Route: Route | undefined;
    public lastUpdate: Date;
    public tracker_id: number;
    public location: Location | null;

    constructor(id: number, name: string, created: Date, updated: Date, enabled: boolean, trackerID: number) {
        this.id = id;
        this.name = name;
        this.created = created;
        this.updated = updated;
        this.enabled = enabled;
        this.lat = 0;
        this.lng = 0;
        this.heading = 0;
        this.speed = 0;
        this.RouteID = null;
        this.marker = new L.Marker([this.lat, this.lng], {
            icon: L.icon({
                iconUrl: getMarkerString('#FFF'),
                iconSize: [32, 32], // size of the icon
                iconAnchor: [16, 16], // point of the icon which will correspond to marker's location
                popupAnchor: [0, 0],   // point from which the popup should open relative to the iconAnchor
            }),
            zIndexOffset: 1000,
            rotationOrigin: 'center',
        });
        this.shownOnMap = false;
        this.map = undefined;
        this.Route = undefined;
        this.lastUpdate = new Date();
        this.tracker_id = trackerID;
        this.location = null;
    }

    public getMessage(): string {
        if (this.location === null) {
            return '';
        }
        const speed = Math.round(this.location.speed * 100) / 100;
        const direction = getCardinalDirection(this.location.heading);
        const routeOnMsg = this.Route === undefined ? '' : `on route <i>${this.Route.name}</i>`;
        let message = `<b>${this.name}</b> ${routeOnMsg}<br>`
            + `Traveling ${direction} at ${speed} mph`;
        if (this.location !== undefined) {
            message += '<br>as of ' + this.location.time.toLocaleTimeString();
        }
        return message;
    }

    public addToMap(map: L.Map) {
        if (this.map === undefined) {
            this.map = map;
        }
        this.marker.bindPopup(this.getMessage());
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

    public setHeading(heading: number) {
        this.heading = heading;
        this.marker.setRotationAngle(this.heading - 45);
        this.marker.bindPopup(this.getMessage());
    }

    public setRoute(r: Route | undefined) {
        console.log('set route');
        if (r === undefined) {
            this.marker.setIcon(L.icon({
                iconUrl: getMarkerString('#FFF'),
                iconSize: [32, 32], // size of the icon
                iconAnchor: [16, 16], // point of the icon which will correspond to marker's location
                popupAnchor: [0, 0],   // point from which the popup should open relative to the iconAnchor
            }));

            return;
        }
        this.Route = r;
        this.RouteID = r.id;
        this.marker.setIcon(L.icon({
            iconUrl: getMarkerString(r.color),
            iconSize: [32, 32], // size of the icon
            iconAnchor: [16, 16], // point of the icon which will correspond to marker's location
            popupAnchor: [0, 0],   // point from which the popup should open relative to the iconAnchor
        }));
        this.marker.bindPopup(this.getMessage());

    }

    public setLatLng(lat: number, lng: number) {
        this.lat = lat;
        this.lng = lng;
        this.marker.setLatLng([this.lat, this.lng]);
    }

    public removeFromMap(map: L.Map) {
        map.removeLayer(this.marker);
    }

    public asJSON(): { id: number; tracker_id: string; name: string; enabled: boolean } {
        return {
            id: this.id,
            enabled: this.enabled,
            tracker_id: String(this.tracker_id),
            name: this.name,
        };
    }

    public setLocation(location: Location) {
        this.location = location;

        // update marker
        this.setLatLng(this.location.latitude, this.location.longitude);
        this.setHeading(this.location.heading);
        this.speed = this.location.speed;

        const now = new Date().getTime();
        const fiveMinMilliseconds = 5 * 60 * 1000;
        if (now - this.location.time.getTime() > fiveMinMilliseconds) {
            this.showOnMap(false);
        } else {
            this.showOnMap(true);
        }
    }
}
