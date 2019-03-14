import Route, { RouteInterface } from '../route';
import Vehicle from '../vehicle';
import { Stop } from '../stop';
import AdminMessageUpdate from '../adminMessageUpdate';

export default class AdminServiceProvider {
    public static EditRoute(route: Route): Promise<Response> {
        return fetch('/routes/edit', {
            method: 'POST',
            body: JSON.stringify(route as RouteInterface),
        });
    }

    public static DeleteRoute(route: Route): Promise<Response> {
        return fetch('/routes?id=' + String(route.id), {
            method: 'DELETE',
        });
    }

    public static CreateRoute(route: Route): Promise<Response> {
        return fetch('/routes/create', {
            method: 'POST',
            body: JSON.stringify(route as RouteInterface),
        });
    }

    public static EditVehicle(vehicle: Vehicle): Promise<Response> {
        return fetch('/vehicles/edit', {
            method: 'POST',
            body: JSON.stringify(vehicle.asJSON()),
        });
    }

    public static DeleteVehicle(vehicle: Vehicle): Promise<Response> {
        return fetch('/vehicles?id=' + String(vehicle.id), {
            method: 'DELETE',
        });
    }

    public static NewVehicle(vehicle: Vehicle): Promise<Response> {
        return fetch('/vehicles/create', {
            method: 'POST',
            body: JSON.stringify(vehicle.asJSON()),
        });
    }

    public static NewStop(stop: Stop): Promise<Response> {
        return fetch('/stops/create', {
            method: 'POST',
            body: JSON.stringify(stop.asJSON()),
        });
    }

    public static SetMessage(message: AdminMessageUpdate): Promise<Response> {
        return fetch('/adminMessage', {
            method: 'POST',
            body: JSON.stringify(message),
        });
    }
}
