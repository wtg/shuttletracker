import Vehicle from '../vehicle';
import Route from '../route';
import Stop from '../stop';
/**
 * Info service provider grabs the basic information from the api and returns it to the frontend.
 */
export default class InfoServiceProvider {
    public GrabVehicles(): Promise<Vehicle[]> {
        return fetch('http://shuttles.rpi.edu/vehicles').then((data) => data.json()).then((data) => {
            const ret = new Array<Vehicle>();
            data.forEach((element: {
                vehicleID: number,
                vehicleName: string,
                Created: string,
                Updated: string,
                enabled: boolean,
            }) => {
                ret.push(new Vehicle(element.vehicleID, element.vehicleName,
                    new Date(element.Created), new Date(element.Updated), element.enabled));
            });
            return ret;
        });

    }

    public GrabRoutes(): Promise<Route[]> {
        return fetch('http://shuttles.rpi.edu/routes').then((data) => data.json()).then((data) => {
            const ret = new Array<Route>();
            data.forEach((element: {
                id: string,
                name: string,
                description: string,
                intervals: any[],
                enabled: boolean,
                active: boolean,
                color: string,
                width: string,
                coords: [{lat: number, lng: number}],
            }) => {
                ret.push(new Route(element.id, element.name, element.description,
                    element.enabled, element.active, element.color, Number(element.width), element.coords));
            });
            return ret;
        });
    }

    public GrabStops(): Promise<Stop[]> {
        return fetch('http://shuttles.rpi.edu/stops').then((data) => data.json()).then((data) => {
            const ret = new Array<Stop>();
            data.forEach((element: {
                id: string,
                name: string,
                description: string,
                lat: string,
                lng: string,
                address: string,
                startTime: string,
                endTime: string,
                enabled: string,
                routeId: string,
                segmentIndex: number,
            }) => {
                ret.push(new Stop(element.id, element.name, element.description, Number(element.lat),
                                    Number(element.lng), element.address, element.startTime, element. endTime,
                                    Boolean(element.enabled), element.routeId, element.segmentIndex));
            });
            return ret;
        });
    }
}
