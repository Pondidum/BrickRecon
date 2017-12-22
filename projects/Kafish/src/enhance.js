import uuid from "uuid";

const replaceStrings = event => {
  Object.keys(event)
    .filter(prop => event[prop] === "")
    .forEach(prop => {
      event[prop] = null;
    });

  return event;
};

export default event =>
  Object.assign(
    { eventId: uuid() },
    replaceStrings(typeof event === "string" ? JSON.parse(event) : event),
    { timestamp: new Date().getTime() }
  );
