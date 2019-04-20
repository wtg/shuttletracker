export default class Location {
    public id: number;
    public vehicleID: number;
    public created: Date;
    public time: Date;
    public latitude: number;
    public longitude: number;
    public heading: number;
    public speed: number;
    public routeID: number;

    constructor(id: number, vehicleID: number, created: Date, time: Date, latitude: number, longitude: number, heading: number, speed: number, routeID: number) {
        this.id = id;
        this.vehicleID = vehicleID;
        this.created = created;
        this.time = time;
        this.latitude = latitude;
        this.longitude = longitude;
        this.heading = heading;
        this.speed = speed;
        this.routeID = routeID;
    }
}
