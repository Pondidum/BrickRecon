import { DynamoDB } from "aws-sdk";
const dynamo = new DynamoDB.DocumentClient();

const write = (tableName, contents) =>
  dynamo.put({ TableName: tableName, Item: contents }).promise();

const read = (tableName, timestamp) =>
  dynamo
    .scan({
      TableName: tableName,
      FilterExpression: "#ts > :ts",
      ExpressionAttributeNames: {
        "#ts": "timestamp"
      },
      ExpressionAttributeValues: {
        ":ts": timestamp
      }
    })
    .promise();

module.exports = tableName => {
  return {
    write: contents => write(tableName, contents),
    read: timestamp => read(tableName, timestamp)
  };
};
