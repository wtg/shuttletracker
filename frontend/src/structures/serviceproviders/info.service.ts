import Vehicle from '../vehicle';
import Route from '../route';
import Stop from '../stop';
import Update from '../update';
import AdminMessageUpdate from '@/structures/adminMessageUpdate';
/**
 * Info service provider grabs the basic information from the api and returns it to the frontend.
 */
export default class InfoServiceProvider {
    public GrabVehicles(): Promise < Vehicle[] > {
        return fetch('https://shuttles.rpi.edu/vehicles').then((data) => data.json()).then((data) => {
            const ret = new Array < Vehicle > ();
            data.forEach((element: {
                id: number,
                name: string,
                created: string,
                updated: string,
                enabled: boolean,
                tracker_id: string,
            }) => {
                ret.push(new Vehicle(element.id, element.name,
                    new Date(element.created), new Date(element.updated), element.enabled));
            });
            return ret;
        });

    }

    public GrabRoutes(): Promise < Route[] > {
        return fetch('https://shuttles.rpi.edu/routes').then((data) => data.json()).then((data) => {
            const ret = new Array < Route > ();
            data.forEach((element: {
                id: number,
                name: string,
                description: string,
                enabled: boolean,
                color: string,
                width: number,
                points: [{
                    latitude: number,
                    longitude: number,
                }],
            }) => {
                ret.push(new Route(element.id, element.name, element.description,
                    element.enabled, element.color, Number(element.width), element.points));
            });
            return ret;
        });
    }

    public GrabStops(): Promise < Stop[] > {
        return fetch('https://shuttles.rpi.edu/stops').then((data) => data.json()).then((data) => {
            const ret = new Array < Stop > ();
            data.forEach((element: {
                id: number,
                name: string,
                description: string,
                latitude: string,
                longitude: string,
                created: string,
                updated: string,
            }) => {
                ret.push(new Stop(element.id, element.name, element.description, Number(element.latitude),
                    Number(element.longitude), element.created, element.updated));
            });
            return ret;
        });
    }

    public GrabAdminMessage(): Promise <AdminMessageUpdate> {
        return fetch('https://shuttles.rpi.edu/adminMessage').then((data) => data.json()).then((ret) => {
            return new AdminMessageUpdate(0, '', ret.message, ret.enabled, new Date(ret.created));
        });
    }

    public GrabUpdates(): Promise < Update[] > {
        return fetch('https://shuttles.rpi.edu/updates').then((data) => data.json()).then((data): Update[] => {
            const ret = new Array <Update> ();
            data.forEach((element: Update) => {
                ret.push(element);
            });
            return ret;
        });
    }
}
