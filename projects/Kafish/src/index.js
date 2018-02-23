import writer from "./writer";
import reader from "./reader";
import storage from "./storage";
import Notifier from "./notifier";

const tableName = process.env.TABLE_NAME;
const snsTopic = process.env.SNS_TOPIC;

const store = storage(tableName);
const notifier = new Notifier({ topic: snsTopic });

const response = (status, body) => {
  return {
    statusCode: status,
    body: JSON.stringify(body)
  };
};

export const writeHandler = (awsEvent, context, callback) =>
  writer({
    store,
    publish: notifier.publish,
    awsEvent,
    respond: (status, body) => callback(null, response(status, body))
  });

export const readHandler = (awsEvent, context, callback) =>
  reader({
    store,
    awsEvent,
    respond: (status, body) => callback(null, response(status, body))
  });
