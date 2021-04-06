import Route from '@/structures/route';
import Vehicle from '@/structures/vehicle';
import { Stop } from '@/structures/stop';
import ETA from '@/structures/eta';
import Form from '@/structures/form';
import AdminMessageUpdate from '@/structures/adminMessageUpdate';

// Define types for the vuex state store
export interface StoreState {
    Routes: Route[];
    Vehicles: Vehicle[];
    Stops: Stop[];
    etas: ETA[];
    Forms: Form[];
    adminMessage: AdminMessageUpdate | undefined;
    online: boolean;
    now: Date;
    settings: {
        busButtonEnabled: boolean,
        etasEnabled: boolean,
        fusionPositionEnabled: boolean,
        busButtonChoice: string,
        darkThemeMode: string,
    };
    geolocationDenied: boolean;

    // This has three states: it is initially undefined, then it gets set to true
    // after Fusion client connects, and it gets set to false if it disconnects.
    fusionConnected: boolean | undefined;
}
