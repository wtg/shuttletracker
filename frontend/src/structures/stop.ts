import * as L from 'leaflet';
import Route from '@/structures/route';
import {StoreState} from '@/StoreState';
import {DarkTheme} from '@/structures/theme';

/**
 * SVGs used for representing a stop on the map and in the legend
 */
export const StopSVGLight = require('@/assets/circle.svg') as string;
export const StopSVGDark = require('@/assets/circle.svg') as string;

/**
 * Stop represents a single stop on a route
 */
export class Stop {

    public static createMarkerIconForCurrentTheme(state: StoreState): L.Icon {
        return L.icon({
            iconUrl: DarkTheme.isDarkThemeVisible(state) ? StopSVGDark : StopSVGLight,
            iconSize: [12, 12], // size of the icon
            iconAnchor: [6, 6], // point of the icon which will correspond to marker's location
            shadowAnchor: [6, 6], // the same for the shadow
            popupAnchor: [0, 0], // point from which the popup should open relative to the iconAnchor
        });
    }

    public id: number;
    public name: string;
    public description: string;
    public latitude: number;
    public longitude: number;
    public created: string;
    public updated: string;
    public routesOn: Route[];
    public marker: L.Marker | null;

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
        this.marker = null;
    }

    public getMessage(): string {
        if (this.routesOn.length > 0) {
            return this.name +
                ` is on route${this.routesOn.length > 1 ? 's' : ''} `
                + this.routesOn.filter((route: Route) => route.shouldShow()).map((route: Route) => `<i>${route.name}</i>`).join(', ');
        } else {
            return this.name;
        }
    }

    public getOrCreateMarker(state: StoreState): L.Marker {
        if (this.marker === null) {
            this.marker = L.marker([this.latitude, this.longitude], {
                icon: Stop.createMarkerIconForCurrentTheme(state),
            });
        }
        return this.marker;
    }

    public addRoute(route: Route): void {
        this.routesOn.push(route);
    }

    public containsRoute(route_id: number): boolean {
        for (const route of this.routesOn) {
            if (route.id === route_id) {
                return true;
            }
        }
        return false;
    }

    public shouldShow(): boolean {
        for (const route of this.routesOn) {
            if (route.shouldShow()) {
                return true;
            }
        }
        return false;
    }

    public asJSON(): {
        id: number, name: string; description: string; latitude: number; longitude: number } {
        return {
            id: this.id,
            name: this.name,
            description: this.description,
            latitude: Number(this.latitude),
            longitude: Number(this.longitude),
        };
    }
}
