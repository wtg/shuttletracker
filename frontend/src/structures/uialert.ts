export enum Color {
    Red,
    Yellow,
    Green,
    Blue,
}

export default class UIAlert {
    public message: string;
    public color: Color;

    constructor(message: string, color: Color) {
        this.color = color;
        this.message = message;
    }
}
