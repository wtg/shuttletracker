const BASE_URL = process.env.BASE_URL as string;
const ROOT_URL = BASE_URL.slice(0, -8);
function constructURL(path: string): string {
    return ROOT_URL + path;
}

export default {
    UpdatesURL: constructURL('updates'),
    StopsURL: constructURL('stops'),
    RoutesURL: constructURL('routes'),
    VehiclesURL: constructURL('vehicles'),
    AdminMessageURL: constructURL('adminMessage'),
    FusionURL: constructURL('fusion/'),
};
