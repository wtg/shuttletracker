export default class AdminMessageUpdate {

    public id: number;
    public Type: string;
    public Message: string;
    public Display: boolean;
    public Created: Date;

    constructor(id: number, Type: string, Message: string, Display: boolean, Created: Date) {
        this.id = id;
        this.Type = Type;
        this.Message = Message;
        this.Display = Display;
        this.Created = Created;

    }

}
