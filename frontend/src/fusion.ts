export default class Fusion {
    constructor() {
        console.log("Fusion created.");
        const ws = new WebSocket("ws://localhost:8080/fusion");
        ws.onopen = (event) => {
            ws.send("hey fusion");
        };

        const options = { enableHighAccuracy: true, maximumAge: 0 };
        navigator.geolocation.watchPosition(
            position => {
                const data = {
                    latitude: position.coords.latitude,
                    longitude: position.coords.longitude,
                    heading: position.coords.heading,
                    speed: position.coords.speed
                };
                ws.send(JSON.stringify(data));
            },
            error => {
                console.log("could not get position", error);
            }, options
        )
    }
}
