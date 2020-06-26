export interface FeedbackInterface {
    id: number;
    topic: string;
    message: string;
    enabled: boolean;
}

/**
 * Form represents a single form
 */
export default class Form {
    public id: number;
    public topic: string;
    public message: string;
    public read: boolean;

    constructor(id: number, topic: string, message: string, read: boolean) {
        this.id = id;
        this.topic = topic;
        this.message = message;
        this.read = read;
    }

    // returns feedback message and marks the form as read
    public getMessage(): string {
        this.read = true;
        return this.message;
    }

    // sets form to unread if needed to
    public unread() {
        this.read = false;
    }

    public asJSON(): { id: number; topic: string; message: string; read: boolean } {
        return {
            id: this.id,
            topic: String(this.topic),
            message: String(this.message),
            read: this.read,
        };
    }

}
