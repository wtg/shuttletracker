import * as L from 'leaflet';
import Route from '@/structures/route';

const StopSVG = require('@/assets/circle.svg') as string;

/**
 * Stop represents a single stop on a route
 */
export default class Stop {
    public id: number;
    public name: string;
    public description: string;
    public latitude: number;
    public longitude: number;
    public created: string;
    public updated: string;
    public routesOn: Route[];
    public marker: L.Marker;

    constructor(id: number, name: string, description: string,
                lat: number, lng: number, created: string, updated: string) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.latitude = lat;
        this.longitude = lng;
        this.created = created;
        this.updated = updated;
        this.routesOn = [];
        this.marker = L.marker([this.latitude, this.longitude], {
            icon: L.icon({
                iconUrl: StopSVG,
                iconSize: [12, 12], // size of the icon
                iconAnchor: [6, 6], // point of the icon which will correspond to marker's location
                shadowAnchor: [6, 6], // the same for the shadow
                popupAnchor: [0, 0], // point from which the popup should open relative to the iconAnchor
              }),
        });
    }

    public getMessage(): string {
        const msg = this.name +
                  ` is${this.routesOn.length > 0 ? ' on' : 'n\'t on a'} route${this.routesOn.length > 1 ? 's' : ''} `
                  + this.routesOn.map((route: Route) => route.name).join(', ');
        return msg;
    }

    public addRoute(route: Route): void {
        this.routesOn.push(route);
    }
}
