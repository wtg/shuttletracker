import Route from '@/structures/route';
import Vehicle from '@/structures/vehicle';
import { Stop } from '@/structures/stop';
import AdminMessageUpdate from '@/structures/adminMessageUpdate';

// Define types for the vuex state store
export interface StoreState {
    Routes: Route[];
    Vehicles: Vehicle[];
    Stops: Stop[];
    adminMessage: AdminMessageUpdate | undefined;
    online: boolean;
    settings: {
        busButtonEnabled: boolean,
        fusionPositionEnabled: boolean,
    };
    geolocationDenied: boolean;
}
