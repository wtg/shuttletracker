const dateToSttr: string[] = ["Sunday","Monday","Tuesday","Wednesday","Thursday","Friday","Saturday"];

// A helper class for schedules, corresponds to routeActiveInternal
export default class routeScheduleInterval {
    public id: number;
    public routeID: number;
    public startTime: Date;
    public startDate: number;
    public endDate: number;
    public endTime: Date;

    constructor(id: number, routeID: number, startDay: number, startTime: Date, endDay: number, endTime: Date) {
        this.id = id;
        this.routeID = routeID;
        this.startDate = startDay;
        this.startTime = startTime;
        this.endDate = endDay;
        this.endTime = endTime;
    }

    // Make the interval into a readable string
    public toString(): string{
        return dateToSttr[this.startDate] + " at " + this.endTime.getHours() + " to " + dateToSttr[this.endDate] + " at " + this.endTime;
    }
}
