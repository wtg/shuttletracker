// An update to a vehicle
export default interface Update {
    id: number;
    latitude: number;
    longitude: number;
    heading: number;
    speed: number;
    time: string;
    date: string;
    status: string;
    created: string;
    route_id: number | null;
    vehicle_id: number;

}
