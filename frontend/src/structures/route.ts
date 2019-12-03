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

        // To get scheduling to work, you need to call shouldActive at regular intervals
        // We were unable to test this because the routes they gave us didn't have schedules
        // attached to them
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

    // checks whether the route should show up using the current time and the time intervals
    // that should be active.
    public shouldActive(): boolean {
        // this should get the current Date & Time
        const currentTime = new Date();
        // this should go through all of the active intervals that the route should be active for
        for (const thisInterval of this.schedule) {
            // this should check if the current time is within any of those intervals
            if (currentTime.getTime() >= thisInterval.start_time.getTime() &&
                currentTime.getTime() <= thisInterval.end_time.getTime()) {
                // setting this to active should make shouldShow return false, which should
                // toggle the route in Public.vue
                this.active = true;
                return true;
            }
        }

        this.active = false;
        return false;
    }

}
