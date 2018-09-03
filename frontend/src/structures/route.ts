/**
 * This represents a single route as returned by the shuttle tracker
 */
export default class Route {
    public id: string;
    public name: string;
    public description: string;
    public enabled: boolean;
    public active: boolean;
    public color: string;
    public width: number;
    public coords: [{
        lat: number,
        lng: number,
    }];
    constructor(id: string, name: string, description: string, enabled: boolean,
                active: boolean, color: string, width: number, coords: [{
                    lat: number,
                    lng: number,
        }]) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.enabled = enabled;
        this.active = active;
        this.color = color;
        this.width = Number(width);
        this.coords = coords;
    }
}
