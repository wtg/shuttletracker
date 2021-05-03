// FeedbackMessageUpdate represents an update from the backend for the feedback message
export default class FeedbackMessageUpdate {

    public message: string;
    public admin: boolean;

    constructor(message: string, admin: boolean) {
        this.message = message;
        this.admin = admin;
    }

}
