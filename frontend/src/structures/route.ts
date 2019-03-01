import routeScheduleInterval from './routeScheduleInterval';

export interface RouteInterface {
    id: number;
    name: string;
    description: string;
    enabled: boolean;
    color: string;
    width: number;
    schedule: routeScheduleInterval[];
    points: Array<{
        latitude: number,
        longitude: number,
    }>;
}

/**
 * This represents a single route as returned by the shuttle tracker
 */
export default class Route implements RouteInterface {
    public id: number;
    public name: string;
    public description: string;
    public enabled: boolean;
    public color: string;
    public width: number;
    public schedule: routeScheduleInterval[];
    public active: boolean;
    public points: Array<{
        latitude: number,
        longitude: number,
    }>;

    constructor(id: number, name: string, description: string, enabled: boolean,
                color: string, width: number, points: Array<{
            latitude: number,
            longitude: number,
        }>,     schedule: routeScheduleInterval[], active: boolean) {

        this.active = active;
        this.id = id;
        this.name = name;
        this.description = description;
        this.enabled = enabled;
        this.color = color;
        this.width = Number(width);
        this.points = points;
        this.schedule = schedule;
    }

    public shouldShow(): boolean {
        return this.enabled && this.active;
    }



}
