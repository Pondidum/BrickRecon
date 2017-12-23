import writer from "./writer";
import reader from "./reader";
import storage from "./storage";
import notification from "./notification";

const tableName = process.env.TABLE_NAME;
const snsTopic = process.env.SNS_TOPIC;

const store = storage(tableName);
const notify = notification(snsTopic);

const response = (status, body) => {
  return {
    statusCode: status,
    body: JSON.stringify(body)
  };
};

export const writeHandler = (awsEvent, context, callback) =>
  writer({
    store,
    publish: notify.publish,
    awsEvent,
    respond: (status, body) => callback(null, response(status, body))
  });

export const readHandler = (awsEvent, context, callback) =>
  reader({
    store,
    publish: notify.publish,
    awsEvent,
    respond: (status, body) => callback(null, response(status, body))
  });
