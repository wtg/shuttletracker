import store from '@/store';

// UserLocationService allows callbacks to be registered on user location updates
export default class UserLocationService {

    public static getInstance(): UserLocationService {
        if (UserLocationService.instance === undefined) {
            UserLocationService.instance = new UserLocationService();
        }
        return UserLocationService.instance;
    }

    private static instance: UserLocationService;
    private callbacks: Array<((pos: Position) => any)>;
    private currentPosition: Position | undefined;

    private constructor() {
        this.callbacks = new Array<((pos: Position) => any)>();
        this.startGeolocation();
    }

    // registerCallback registers a function to be called on all location updates
    public registerCallback(callback: (pos: Position) => any) {
        this.callbacks.push(callback);
    }

    public getCurrentLocation(): Position | undefined {
        return this.currentPosition;
    }

    private startGeolocation() {
        if (!('geolocation' in navigator)) {
            console.log('Client does not support geolocation.');
            return;
        }
        const options = { enableHighAccuracy: true, maximumAge: 0 };
        navigator.geolocation.watchPosition(
            (position: Position) => {
                this.currentPosition = position;
                for (let i = 0; i < this.callbacks.length; i++) {
                    this.callbacks[i](position);
                }
            },
            (error) => {
                if (error.code === 1) {
                    store.commit('setGeolocationDenied', true);
                }
            }, options,
        );
    }

}
