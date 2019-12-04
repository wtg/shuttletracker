// AdminMessageUpdate represents an update from the backend for the admin message
export default class AnnouncementUpdate {

  public message: string;
  public enabled: boolean;
  public created: Date;
  public updated: Date;
  public link: string;

  constructor(message: string, enabled: boolean, created: Date, updated: Date, link: string) {
      this.message = message;
      this.enabled = enabled;
      this.created = created;
      this.updated = updated;
      this.link = link;

  }

}
