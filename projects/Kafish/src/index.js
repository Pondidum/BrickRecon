import writer from "./writer";
import reader from "./reader";
import storage from "./storage";

const tableName = process.env.TABLE_NAME;
const store = storage(tableName);

const response = (status, body) => {
  return {
    statusCode: status,
    body: JSON.stringify(body)
  };
};

export const writeHandler = (awsEvent, context, callback) =>
  writer({
    store,
    awsEvent,
    respond: (status, body) => callback(null, response(status, body))
  });

export const readHandler = (awsEvent, context, callback) =>
  reader({
    store,
    awsEvent,
    respond: (status, body) => callback(null, response(status, body))
  });
