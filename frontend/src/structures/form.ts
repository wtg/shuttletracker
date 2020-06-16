export interface FeedbackInterface {
    id: number;
    message: string;
    enabled: boolean;
}

/**
 * Form represents a single form
 */
export default class Form {
    public id: number;
    public message: string;
    public read: boolean;

    constructor(id: number, message: string, read: boolean) {
        this.id = id;
        this.message = message;
        this.read = read;
    }

    // returns feedback message and marks the form as read
    public getMessage(): string {
        this.read = true;
        return this.message;
    }

    //sets form to unread if needed to
    public unread() {
        this.read = false;
    }

    public asJSON(): { id: number; message: string; read: boolean } {
        return {
            id: this.id,
            message: String(this.message),
            read: this.read,
        };
    }

}