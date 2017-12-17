import uuid from "uuid";
import { DynamoDB } from "aws-sdk";

const dynamo = new DynamoDB.DocumentClient();
const tableName = process.env.TABLE_NAME;

const enhance = event =>
  Object.assign(
    { eventId: uuid() },
    typeof event === "string" ? JSON.parse(event) : event,
    { timestamp: new Date().getTime() }
  );

const write = event =>
  dynamo
    .put({
      TableName: tableName,
      Item: event
    })
    .promise();

export const handler = (awsEvent, context, callback) => {
  const event = enhance(awsEvent.body);

  return write(event)
    .then(() => callback(null, { statusCode: "200", body: "{}" }))
    .catch(err => {
      console.log("error writing to dynamo", JSON.stringify(err, null, 2));
      callback(null, {
        statusCode: "400",
        body: JSON.stringify({
          message: "Unable to store event",
          exception: err
        })
      });
    });
};
