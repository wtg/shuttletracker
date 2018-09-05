import Route from './structures/route';
import Vehicle from './structures/vehicle';
import Stop from '@/structures/stop';

// Define types for the vuex state store
export interface StoreState {
    Routes: Route[];
    Vehicles: Vehicle[];
    Stops: Stop[];
}
