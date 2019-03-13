// import Stop from '@/structures/stop';
// import Vehicle from '@/structures/vehicle';

export default class ETA {
    // public stop: Stop;
    // public vehicle: Vehicle;
    public stopID: number;
    public vehicleID: number;
    public routeID: number;
    public eta: Date;

    constructor(stop: number, vehicle: number, eta: Date, routeID: number) {
        this.stopID = stop;
        this.vehicleID = vehicle;
        this.eta = eta;
        this.routeID = routeID;
    }
}
