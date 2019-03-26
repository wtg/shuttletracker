import Route from '@/structures/route';
import Vehicle from '@/structures/vehicle';
import { Stop } from '@/structures/stop';
import ETA from '@/structures/eta';
import AdminMessageUpdate from '@/structures/adminMessageUpdate';

// Define types for the vuex state store
export interface StoreState {
    Routes: Route[];
    Vehicles: Vehicle[];
    Stops: Stop[];
    etas: ETA[];
    adminMessage: AdminMessageUpdate | undefined;
    online: boolean;
    settings: {
        busButtonEnabled: boolean,
        etasEnabled: boolean,
        fusionPositionEnabled: boolean,
    };
    geolocationDenied: boolean;
}
