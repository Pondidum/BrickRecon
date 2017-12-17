import uuid from "uuid";
import { DynamoDB } from "aws-sdk";

const dynamo = new DynamoDB.DocumentClient();

const enhance = event =>
  Object.assign(
    { eventId: uuid() },
    typeof event === "string" ? JSON.parse(event) : event,
    { timestamp: new Date().getTime() }
  );

const write = (tableName, event) =>
  dynamo.put({ TableName: tableName, Item: event }).promise();

const error = (message, err) =>
  console.log(message, JSON.stringify(err, null, 2));

export default options =>
  write(options.tableName, enhance(options.awsEvent.body))
    .then(() => options.respond("200", {}))
    .catch(err => {
      error("error writing to dynamo", err);

      options.respond("400", {
        message: "Unable to store event",
        exception: err
      });
    });
