import RouteScheduleInterval from './routeScheduleInterval';

export interface RouteInterface {
    id: number;
    name: string;
    description: string;
    enabled: boolean;
    color: string;
    width: number;
    schedule: RouteScheduleInterval[];
    points: Array<{
        latitude: number,
        longitude: number,
    }>;
    stop_ids: number[];
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
    public schedule: RouteScheduleInterval[];
    public active: boolean;
    public points: Array<{
        latitude: number,
        longitude: number,
    }>;
    public stop_ids: number[];

    constructor(id: number, name: string, description: string, enabled: boolean,
                color: string, width: number, points: Array<{
            latitude: number,
            longitude: number,
        }>,     schedule: RouteScheduleInterval[], active: boolean, stop_ids: number[]) {

        this.active = active;
        this.id = id;
        this.name = name;
        this.description = description;
        this.enabled = enabled;
        this.color = color;
        this.width = Number(width);
        this.points = points;
        this.schedule = schedule;
        this.stop_ids = stop_ids;
    }

    public shouldShow(): boolean {
        return this.enabled && this.active;
    }

    public containsStop(stop_id: number): boolean {
        for (const id of this.stop_ids) {
            if (id === stop_id) {
                return true;
            }
        }
        return false;
    }

    public containsTime(): boolean {
        const currentTime = new Date();
        for (const thisInterval of this.schedule) {
            if (currentTime.getTime() >= thisInterval.start_time.getTime() && currentTime.getTime() <= thisInterval.end_time.getTime()) {
                this.active = true;
                return true;
            }
        }
        this.active = false;
        return false;
    }

}
