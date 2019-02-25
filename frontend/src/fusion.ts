import UserLocationService from '@/structures/userlocation.service';
import store from '@/store';

// SocketManager wraps a WebSocket in order to provide guarantees about
// reliability, reconnections, retries, etc.
class SocketManager {
    private ws: WebSocket | null = null;
    private url: string;
    private callbacks = new Array<(data: any) => any>();

    constructor(url: string) {
        this.url = url;
    }

    public registerMessageReceivedCallback(callback: (data: any) => any) {
        this.callbacks.push(callback);
    }

    public open() {
        this.createSocket().then((ws) => {
            this.ws = ws;
        });
    }

    public send(msg: string) {
        if (this.ws) {
            this.ws.send(msg);
        }
    }

    private createSocket(): Promise<WebSocket> {
        return new Promise((resolve, reject) => {
            const ws = new WebSocket(this.url);
            ws.onopen = (event) => {
                // console.log("socket connected", event);
                resolve(ws);
            };
            ws.onmessage = (event) => {
                for (const callback of this.callbacks) {
                    callback(event.data);
                }
            };
            ws.onerror = (event) => {
                // console.log("socket error", event);
            };
            ws.onclose = (event) => {
                // console.log("socket closed", event);
                this.open();
            };
        });
    }
}

export default class Fusion {
    public ws: SocketManager;
    public track = this.generateUUID();
    private callbacks = Array<(message: {}) => any>();

    constructor() {
        const wsURL = this.relativeWSURL('fusion/');
        this.ws = new SocketManager(wsURL);
        this.ws.registerMessageReceivedCallback((data) => {
            const message = JSON.parse(data);
            for (const callback of this.callbacks) {
                callback(message);
            }
        });
    }

    public start() {
        this.ws.open();
        // register location callback
        UserLocationService.getInstance().registerCallback((position) => {
            if (!store.state.settings.fusionPositionEnabled) {
                return;
            }
            const data = {
                type: 'position',
                message: {
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude,
                    heading: position.coords.heading,
                    speed: position.coords.speed,
                    track: this.track,
                },
            };
            this.ws.send(JSON.stringify(data));
        });
    }

    public registerMessageReceivedCallback(callback: (message: {}) => any) {
        this.callbacks.push(callback);
    }

    public sendBusButton() {
        const ls = UserLocationService.getInstance();
        const pos = ls.getCurrentLocation();
        if (!pos) {
            // client geolocation isn't enabled or known
            return;
        }
        const data = {
            type: 'bus_button',
            message: {
                latitude: pos.coords.latitude,
                longitude: pos.coords.longitude,
            },
        };
        this.ws.send(JSON.stringify(data));
    }

    private generateUUID() {
        let d = new Date().getTime();
        if (typeof performance !== 'undefined' && typeof performance.now === 'function') {
            d += performance.now(); // use high-precision timer if available
        }
        return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
            const r = (d + Math.random() * 16) % 16 | 0;
            d = Math.floor(d / 16);
            return (c === 'x' ? r : (r & 0x3 | 0x8)).toString(16);
        });
    }

    private relativeWSURL(wsURL: string) {
        let url = '';
        if (window.location.protocol === 'https:') {
            url += 'wss:';
        } else {
            url += 'ws:';
        }
        url += '//' + window.location.host + window.location.pathname;
        return url + wsURL;
    }
}
