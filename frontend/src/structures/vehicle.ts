/**
 * Vehicle represents a returned vehicle value
 */
export default class Vehicle {
    public id: number;
    public name: string;
    public created: Date;
    public updated: Date;
    public enabled: boolean;

    constructor(id: number, name: string, created: Date, updated: Date, enabled: boolean) {
        this.id = id;
        this.name = name;
        this.created = created;
        this.updated = updated;
        this.enabled = enabled;
    }
}
