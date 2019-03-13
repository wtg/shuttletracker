// import ETA from '@/structures/eta';

/**
 * Stop represents a single stop on a route
 */
export default class Stop {
    public id: number;
    public name: string;
    public description: string;
    public latitude: number;
    public longitude: number;
    public created: string;
    public updated: string;
    // public etas: ETA[];

    constructor(id: number, name: string, description: string,
                lat: number, lng: number, created: string, updated: string) {
        this.id = id;
        this.name = name;
        this.description = description;
        this.latitude = lat;
        this.longitude = lng;
        this.created = created;
        this.updated = updated;
        // this.etas = [];
    }
}
