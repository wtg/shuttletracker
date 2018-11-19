import Route, { RouteInterface } from '../route';

export default class AdminServiceProvider {
    public static EditRoute(route: Route): Promise<Response> {
        return fetch('/routes/edit', {
            method: 'POST',
            body: JSON.stringify(route as RouteInterface),
        });
    }
}
