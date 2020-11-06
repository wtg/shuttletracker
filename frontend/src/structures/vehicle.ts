import * as L from 'leaflet';
import Route from '@/structures/route';
import Location from '@/structures/location';
import getMarkerString from '@/structures/leaflet/rotatedMarker';
import getCardinalDirection from '@/structures/cardinalDirection';
import 'leaflet-rotatedmarker';
import 'leaflet.marker.slideto';
import store from '@/store';

// vehicles are hidden when their most recent location update becomes this old
const tinycolor = require('tinycolor2');
const vehicleInactiveDurationMS = 5 * 60 * 1000;  // five minutes in milliseconds


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
    public RouteID: number | null;
    public shownOnMap: boolean;
    public map: L.Map | undefined;
    public Route: Route | undefined;
    public lastUpdate: Date;
    public tracker_id: number;
    public location: Location | null;
    private hideTimer: number | null = null;
    private pointIndex: number | null;
    private endPointIndex: number | null;
    private destinationLat: number;
    private destinationLng: number;
    private destinationHeading: number;
    private totalTransitionDistance: number;

    constructor(id: number, name: string, created: Date, updated: Date, enabled: boolean, trackerID: number) {
        this.id = id;
        this.name = name;
        this.created = created;
        this.updated = updated;
        this.enabled = enabled;
        this.lat = 0;
        this.lng = 0;
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
        this.pointIndex = null;
        this.endPointIndex = null;
        this.destinationLat = 0;
        this.destinationLng = 0;
        this.destinationHeading = 0;
        this.totalTransitionDistance = 0;
        (this.marker as any).on('moveend', this.continueMoving, this);
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
        this.marker.setRotationAngle(heading - 45);
        this.marker.bindPopup(this.getMessage());
        this.destinationHeading = heading - 45;
    }

    public setRoute(r: Route | undefined, darkEnabled: boolean) {
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
        let markerColor = r.color;
        if (darkEnabled) {
            const darkColor = tinycolor(r.color);
            darkColor.darken(15);
            markerColor = darkColor.toString();
        }
        this.marker.setIcon(L.icon({
            iconUrl: getMarkerString(r.color),
            iconSize: [32, 32], // size of the icon
            iconAnchor: [16, 16], // point of the icon which will correspond to marker's location
            popupAnchor: [0, 0],   // point from which the popup should open relative to the iconAnchor
        }));
        this.marker.bindPopup(this.getMessage());

    }

    // Immediately moves the shuttle to the given position, without a sliding animation
    public setLatLngImmediate(lat: number, lng: number) {
        this.lat = lat;
        this.lng = lng;
        this.marker.setLatLng([this.lat, this.lng]);
    }

    // Moves the shuttle to the given position over the next 5 seconds, with a sliding animation
    public setLatLng(lat: number, lng: number) {
        if (this.lat === 0 || this.lng === 0) {
            this.lat = lat;
            this.lng = lng;
            this.marker.setLatLng(L.latLng(lat, lng));
            return;
        }
        if (store.state.settings.shuttleSlideEnabled) {
            console.log(this.id + ': Starting');
            this.destinationLat = lat;
            this.destinationLng = lng;
            if (this.pointIndex === null) {
                this.updateClosestStartingPoint();
            }
            this.updateClosestDestinationPoint();
            console.log(this.id + ': Destination lat/lng: ' + lat + ' ' + lng);
            console.log(this.id + ': Start point index: ' + this.pointIndex);
            console.log(this.id + ': End point index: ' + this.endPointIndex);
            if (this.Route !== undefined && this.pointIndex !== null  && this.endPointIndex !== null && this.marker.getLatLng().distanceTo(L.latLng(this.Route.points[this.pointIndex].latitude, this.Route.points[this.pointIndex].longitude)) < 100) {
                if (this.pointIndex !== this.endPointIndex) {
                    this.pointIndex++;
                    if (this.pointIndex >= this.Route.points.length) {
                        this.pointIndex = 0;
                    }
                }
                this.continueMoving();
            } else { // Skip the animation if the shuttle has no route, no closest/destination point, or is over 100 meters away from the route
                this.setLatLngImmediate(lat, lng);
            }
        } else {
            this.setLatLngImmediate(lat, lng);
        }
    }

    // Called when one segment of the sliding animation has completed to begin sliding the shuttle to either
    // the next segment or destination point
    public continueMoving() {
        console.log(this.id + ': Continuing: pointIndex: ' + this.pointIndex);
        if (this.Route !== undefined && this.pointIndex !== null && this.endPointIndex !== null) {
            const origPointIndex = this.pointIndex;
            if (this.pointIndex !== this.endPointIndex) {
                // Shuttle can transition to a route point closer to the destination point
                const point = this.Route.points[this.pointIndex];
                this.lat = point.latitude;
                this.lng = point.longitude;
                this.pointIndex++;
                if (this.pointIndex >= this.Route.points.length) {
                    this.pointIndex = 0;
                }
            } else {
                // Shuttle is as close to the destination point as possible, begin sliding to destination point
                this.lat = this.destinationLat;
                this.lng = this.destinationLng;
                this.endPointIndex = null;
            }
            if (!this.marker.getLatLng().equals(L.latLng(this.lat, this.lng))) {
                // Calculate the time this segment should take, proportional to the distance to be travelled
                const transitionDistance = this.marker.getLatLng().distanceTo(L.latLng(this.lat, this.lng));
                const transitionRatio = transitionDistance / this.totalTransitionDistance;
                const transitionTime = Math.round(transitionRatio * 3000); // 3 s (3000 ms) for the total transition
                // Update shuttle's rotation based on the angle between route points
                const newAngle = this.angleBetween(this.marker.getLatLng().lat, this.marker.getLatLng().lng, this.Route.points[origPointIndex].latitude, this.Route.points[origPointIndex].longitude);
                this.marker.setRotationAngle(newAngle);
                // Begin the next segment
                console.log(this.id + ': Moving to ' + this.lat + ' ' + this.lng);
                (this.marker as any).slideTo([this.lat, this.lng], { duration: transitionTime, keepAtCenter: false });
            } else {
                this.continueMoving();
            }
        } else {
            // Animation has completed
            this.marker.setRotationAngle(this.destinationHeading);
            console.log(this.id + ': Stopping at ' + this.marker.getLatLng().lat + ' ' + this.marker.getLatLng().lng);
        }
    }

    public updateClosestStartingPoint() {
        if (this.Route !== undefined) {
            // Find closest point index to the vehicle's current position
            let minDistance = -1.0;
            for (let i = 0; i < this.Route.points.length; i++) {
                const distance = this.marker.getLatLng().distanceTo(L.latLng(this.Route.points[i].latitude, this.Route.points[i].longitude));
                if (distance < minDistance || minDistance < 0) {
                    minDistance = distance;
                    this.pointIndex = i;
                }
            }
        }
    }

    public updateClosestDestinationPoint() {
        if (this.Route !== undefined) {
            if (this.pointIndex !== null) {
                // Find the closest point index to the vehicle's destination
                this.totalTransitionDistance = 0.0;
                let currentPointIndex = this.pointIndex;
                let nextPointIndex = this.pointIndex + 1;
                if (nextPointIndex >= this.Route.points.length) {
                    nextPointIndex = 0;
                }
                let distanceToNextPoint = this.marker.getLatLng().distanceTo(L.latLng(this.Route.points[currentPointIndex].latitude, this.Route.points[currentPointIndex].longitude));
                while (distanceToNextPoint < L.latLng(this.Route.points[currentPointIndex].latitude, this.Route.points[currentPointIndex].longitude).distanceTo(L.latLng(this.destinationLat, this.destinationLng))) {
                    this.totalTransitionDistance += distanceToNextPoint;
                    currentPointIndex = nextPointIndex;
                    nextPointIndex++;
                    if (nextPointIndex >= this.Route.points.length) {
                        nextPointIndex = 0;
                    }
                    distanceToNextPoint = L.latLng(this.Route.points[currentPointIndex].latitude, this.Route.points[currentPointIndex].longitude).distanceTo(L.latLng(this.Route.points[nextPointIndex].latitude, this.Route.points[nextPointIndex].longitude));
                    if (currentPointIndex === this.pointIndex) { // We've wrapped around to the original point
                        this.totalTransitionDistance = 0.0;
                        break;
                    }
                }
                this.endPointIndex = currentPointIndex;
                this.totalTransitionDistance += L.latLng(this.Route.points[this.endPointIndex].latitude, this.Route.points[this.endPointIndex].longitude).distanceTo(L.latLng(this.destinationLat, this.destinationLng));
            }
        }
    }

    public angleBetween(lat1: number, lng1: number, lat2: number, lng2: number) {
        const radLat1 = this.toRadians(lat1);
        const radLng1 = this.toRadians(lng1);
        const radLat2 = this.toRadians(lat2);
        const radLng2 = this.toRadians(lng2);

        const deltaLongitude = (radLng2 - radLng1);
        const y = Math.sin(deltaLongitude) * Math.cos(radLat2);
        const x = Math.cos(radLat1) * Math.sin(radLat2) - Math.sin(radLat1) * Math.cos(radLat2) * Math.cos(deltaLongitude);

        let brng = Math.atan2(y, x);
        brng = this.toDegrees(brng);
        brng = (brng + 360) % 360;

        return brng - 45;
    }

    public toDegrees(angle: number) {
        return angle * (180 / Math.PI);
    }

    public toRadians(angle: number) {
        return angle * (Math.PI / 180);
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

        this.updateShowOnMap();

        // vehicle hides itself after five min since most recent update
        if (this.hideTimer !== null) {
            window.clearInterval(this.hideTimer);
        }
        const now = new Date().getTime();
        this.hideTimer = window.setTimeout(() => { this.updateShowOnMap(); }, vehicleInactiveDurationMS - (now - location.time.getTime()));
    }

    // hides vehicle if the time of its most recent update is older than five minutes
    public updateShowOnMap() {
        if (this.location === null) {
            this.showOnMap(false);
            return;
        }

        const now = new Date().getTime();
        if (now - this.location.time.getTime() >= vehicleInactiveDurationMS) {
            this.showOnMap(false);
        } else {
            this.showOnMap(true);
        }
    }
}
