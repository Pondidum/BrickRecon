import { DynamoDB } from "aws-sdk";
import writer from "./writer";

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
