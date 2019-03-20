export default class ETA {
    public stopID: number;
    public vehicleID: number;
    public routeID: number;
    public eta: Date;
    public arriving: boolean;

    constructor(stopID: number, vehicleID: number, routeID: number, eta: Date, arriving: boolean) {
        this.stopID = stopID;
        this.vehicleID = vehicleID;
        this.routeID = routeID;
        this.eta = eta;
        this.arriving = arriving;
    }
}
