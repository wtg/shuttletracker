import UserLocationService from '@/structures/userlocation.service';

// SocketManager wraps a WebSocket in order to provide guarantees about
// reliability, reconnections, retries, etc.
class SocketManager {
    public ws: WebSocket | null = null;
    public url: string;

    constructor(url: string) {
        this.url = url;
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
                // console.log(event);
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

    constructor() {
        const wsURL = this.relativeWSURL('fusion/');
        this.ws = new SocketManager(wsURL);
    }

    public start() {
        this.ws.open();
        // register location callback
        UserLocationService.getInstance().registerCallback((position) => {
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
