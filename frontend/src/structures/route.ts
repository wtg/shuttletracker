/**
 * This represents a single route as returned by the shuttle tracker
 */
export default class Route {
    public id: number;
    public name: string;
    public description: string;
    public enabled: boolean;
    public color: string;
    public width: number;
    public active: boolean;
    public coords: [{
        latitude: number,
        longitude: number,
    }];

    constructor(id: number, name: string, description: string, enabled: boolean,
        color: string, width: number, coords: [{
                    latitude: number,
                    longitude: number,
        }], active: boolean) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.enabled = enabled;
        this.color = color;
        this.width = Number(width);
        this.coords = coords;
        this.active = active;
    }

    // shouldShow returns true if the route is active and enabled
    public shouldShow(): boolean{
        return this.active && this.enabled;
    }
}
