import uuid from "uuid";

export default event =>
  Object.assign(
    { eventId: uuid() },
    typeof event === "string" ? JSON.parse(event) : event,
    { timestamp: new Date().getTime() }
  );
