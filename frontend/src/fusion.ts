import UserLocationService from '@/structures/userlocation.service';
import store from '@/store';
import ETA from '@/structures/eta';
import Location from '@/structures/location';

// SocketManager wraps a WebSocket in order to provide guarantees about
// reliability, reconnections, retries, etc.
class SocketManager {
    private ws: WebSocket | null = null;
    private url: string;
    private callbacks = new Array<(data: any) => void>();
    private queue = new Array<string>();
    private reconnectCallbacks = new Array<() => void>();

    constructor(url: string) {
        this.url = url;
    }

    public registerMessageReceivedCallback(callback: (data: any) => void) {
        this.callbacks.push(callback);
    }

    public registerReconnectCallback(callback: () => void) {
        this.reconnectCallbacks.push(callback);
    }

    public open() {
        this.createSocket().then((ws) => {
            this.ws = ws;
            this.flushQueue();
        });
    }

    // Queue a message to send to the WebSocket server, then try to flush the queue.
    public send(msg: string) {
        this.queue.push(msg);
        this.flushQueue();
    }

    // Try to send any messages in the queue to the WebSocket server and delete those
    // which we are able to send.
    private flushQueue() {
        if (!this.ws) {
            return;
        }

        // Keep track of what messages we send so we can remove them from the queue later.
        const sent = new Array<number>();
        for (let i = 0; i < this.queue.length; i++) {
            if (this.ws.readyState !== 1) {
                break;
            }
            this.ws.send(this.queue[i]);
            sent.push(i);
        }

        // Remove sent messages from the queue.
        // Do this in reverse to avoid worrying about messing up indices.
        for (let i = sent.length - 1; i >= 0; i--) {
            this.queue.splice(sent[i], 1);
        }
    }

    private createSocket(): Promise<WebSocket> {
        return new Promise((resolve, reject) => {
            const ws = new WebSocket(this.url);
            ws.onopen = (event) => {
                // console.log("socket connected", event);
                store.commit('setFusionConnected', true);
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

                store.commit('setFusionConnected', false);

                // try to reconnect after a second
                setTimeout(() => {
                    this.createSocket().then((newWS) => {
                        this.ws = newWS;

                        // try to send anything that was queued while the socket was closed
                        this.flushQueue();

                        for (const callback of this.reconnectCallbacks) {
                            callback();
                        }
                    });
                }, 1000);
            };
        });
    }
}

export default class Fusion {
    public ws: SocketManager;
    public track = this.generateUUID();
    private callbacks = Array<(message: {}) => any>();
    private subscriptionTopics = new Set<string>();
    private serverID = null;

    constructor() {
        const wsURL = this.relativeWSURL('fusion/');
        this.ws = new SocketManager(wsURL);
        this.ws.registerMessageReceivedCallback((data) => {
            const message = JSON.parse(data);
            for (const callback of this.callbacks) {
                callback(message);
            }
        });
        this.ws.registerReconnectCallback(() => {
            for (const topic of this.subscriptionTopics) {
                this.requestSubscription(topic);
            }
        });
    }

    public start() {
        this.ws.open();

        // register server ID changed refresh callback
        this.registerMessageReceivedCallback((message: any) => {
            if (message.type !== 'server_id') {
                return;
            }

            // store this server ID or check if it has changed
            if (this.serverID === null) {
                this.serverID = message.message;
            } else if (this.serverID !== message.message) {
                // reload page after a random amount of time
                const wait = Math.random() * 20;
                console.log(`Server ID has changed; reloading after ${wait} seconds`);

                setTimeout(() => {
                    location.reload(true);
                }, wait * 1000);
            }
        });

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

        // subscribe to vehicle location updates
        this.subscribe('vehicle_location');
        this.registerMessageReceivedCallback(this.handleVehicleLocations);

        // get notified of bus button setting changes so we can subscribe to the topic
        store.watch((state) => state.settings.busButtonEnabled, (newValue, oldValue) => {
            if (newValue === true) {
                this.subscribe('bus_button');
            } else {
                this.unsubscribe('bus_button');
            }
        });
        if (store.state.settings.busButtonEnabled === true) {
            this.subscribe('bus_button');
        }

        // subscribe to estimated times of arrival
        store.watch((state) => state.settings.etasEnabled, (newValue, oldValue) => {
            if (newValue === true) {
                this.subscribe('eta');
            } else {
                this.unsubscribe('eta');
            }
        });
        if (store.state.settings.etasEnabled === true) {
            this.subscribe('eta');
        }
        this.registerMessageReceivedCallback(this.handleETAs);
    }

    public registerMessageReceivedCallback(callback: (message: {}) => any) {
        this.callbacks.push(callback);
    }

    public sendBusButton() {
        const ls = UserLocationService.getInstance();
        const pos = ls.getCurrentLocation();
        const emoji = store.state.settings.busButtonChoice;
        if (!pos) {
            // client geolocation isn't enabled or known
            return;
        }
        const data = {
            type: 'bus_button',
            message: {
                latitude: pos.coords.latitude,
                longitude: pos.coords.longitude,
                emojiChoice: emoji,
            },
        };
        this.ws.send(JSON.stringify(data));
    }

    public subscribe(topic: string) {
        this.subscriptionTopics.add(topic);
        this.requestSubscription(topic);
    }

    public unsubscribe(topic: string) {
        this.subscriptionTopics.delete(topic);
        this.requestUnsubscription(topic);
    }

    private requestSubscription(topic: string) {
        const data = {
            type: 'subscribe',
            message: { topic },
        };
        this.ws.send(JSON.stringify(data));
    }

    private requestUnsubscription(topic: string) {
        const data = {
            type: 'unsubscribe',
            message: { topic },
        };
        this.ws.send(JSON.stringify(data));
    }

    private handleETAs(message: any) {
        if (message.type !== 'eta') {
            return;
        }

        const etas = new Array<ETA>();
        for (const stopETA of message.message.stop_etas) {
            const eta = new ETA(
                stopETA.stop_id,
                message.message.vehicle_id,
                message.message.route_id,
                new Date(stopETA.eta),
                stopETA.arriving,
            );
            etas.push(eta);
        }
        store.commit('updateETAs', { vehicleID: message.message.vehicle_id, etas });
    }

    private handleVehicleLocations(message: any) {
        if (message.type !== 'vehicle_location') {
            return;
        }

        const m = message.message;
        const location = new Location(
            m.id,
            m.vehicle_id,
            new Date(m.created),
            new Date(m.time),
            m.latitude,
            m.longitude,
            m.heading,
            m.speed,
            m.route_id,
        );
        store.commit('updateVehicleLocation', location);
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
