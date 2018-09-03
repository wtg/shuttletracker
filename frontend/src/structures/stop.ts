/**
 * Stop represents a single stop on a route
 */
export default class Stop{
    public id: string;
    public name: string;
    public description: string;
    public lat: number;
    public lng: number;
    public address: string;
    public startTime: string;
    public endTime: string;
    public enabled: boolean;
    public routeID: string;
    public segmentindex: number;

    constructor(id: string, name: string, description: string,
                lat: number, lng: number, address: string, startTime: string,
                endTime: string, enabled: boolean, routeID: string, segmentindex: number) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.lat = lat;
        this.lng = lng;
        this.address = address;
        this.startTime = startTime;
        this.endTime = endTime;
        this.enabled = enabled;
        this.routeID = routeID;
        this.segmentindex = segmentindex;
    }
}
