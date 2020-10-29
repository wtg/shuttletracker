// FeedbackMessageUpdate represents an update from the backend for the feedback message
export default class FeedbackMessageUpdate {

    public message: string;
    public enabled: boolean;
    public created: Date;
    public updated: Date;

    constructor(message: string, enabled: boolean, created: Date, updated: Date) {
        this.message = message;
        this.enabled = enabled;
        this.created = created;
        this.updated = updated;

    }

}
