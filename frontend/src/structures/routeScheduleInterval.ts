const dateToSttr: string[] = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday'];

// A helper class for schedules, corresponds to routeActiveInternal
export default class RotueScheduleInterval {
    public id: number;
    public route_id: number;
    public start_time: Date;
    public start_day: number;
    public end_day: number;
    public end_time: Date;

    constructor(id: number, route_id: number, start_day: number, start_time: Date, end_day: number, end_time: Date) {
        this.id = id;
        this.route_id = route_id;
        this.start_day = start_day;
        this.start_time = start_time;
        this.end_day = end_day;
        this.end_time = end_time;
    }

    // Make the interval into a readable string
    public toString(): string {
        return dateToSttr[this.start_day] + ' at ' + this.end_time.getHours() + ':' + this.end_time.getMinutes() + ' to ' + dateToSttr[this.end_day] + ' at ' + this.end_time;
    }
}
