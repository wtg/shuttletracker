export interface FeedbackInterface {
    id: number;
    message: string;
    admin: boolean;
}

/**
 * Form represents a single form
 */
export default class Form {
    public id: number;
    public message: string;
    public created: Date;
    public admin: boolean;

    constructor(id: number, message: string, created: Date, admin: boolean) {
        this.id = id;
        this.message = message;
        this.created = created;
        this.admin = admin;
    }

    // returns form id
    public getID(): number {
        return this.id;
    }

    // returns feedback message
    public getMessage(): string {
        return this.message;
    }

    public when(): Date {
        return this.created;
    }

    public isAdmin(): boolean {
        return this.admin;
    }

    public asJSON(): { message: string; admin: boolean } {
        return {
            message: String(this.message),
            admin: this.admin,
        };
    }

}
