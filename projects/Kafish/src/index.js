import { DynamoDB } from "aws-sdk";
import writer from "./writer";
import reader from "./reader";

const dynamo = new DynamoDB.DocumentClient();
const tableName = process.env.TABLE_NAME;

const response = (status, body) => {
  return {
    statusCode: status,
    body: JSON.stringify(body)
  };
};

export const writeHandler = (awsEvent, context, callback) =>
  writer({
    dynamo,
    tableName: tableName,
    awsEvent,
    respond: (status, body) => callback(null, response(status, body))
  });

export const readHandler = (awsEvent, context, callback) =>
  reader({
    dynamo,
    tableName: tableName,
    awsEvent,
    respond: (status, body) => callback(null, response(status, body))
  });
