import Route, { RouteInterface } from '../route';
import Vehicle from '../vehicle';

export default class AdminServiceProvider {
    public static EditRoute(route: Route): Promise<Response> {
        return fetch('/routes/edit', {
            method: 'POST',
            body: JSON.stringify(route as RouteInterface),
        });
    }

    public static EditVehicle(vehicle: Vehicle): Promise<Response> {
        return fetch('/vehicles/edit', {
            method: 'POST',
            body: JSON.stringify(vehicle.asJSON()),
        })
    }

    public static DeleteVehicle(vehicle: Vehicle): Promise<Response> {
        return fetch('/vehicles?id=' + String(vehicle.id), {
            method: 'DELETE',
        })
    }

    public static NewVehicle(vehicle: Vehicle): Promise<Response> {
        return fetch('/vehicles/create', {
            method: 'POST',
            body: JSON.stringify(vehicle.asJSON()),
        })
    }
}
