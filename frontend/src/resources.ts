// Resources contains methods for retrieving paths to resources on the server.
// This was created to ensure that Shuttle Tracker works regardless of what
// subdirectory it is deployed at.

const BASE_URL = process.env.BASE_URL as string;

export default {
    BasePath: BASE_URL.slice(0, -7), // strip /static/ from Vue's publicPath
};
