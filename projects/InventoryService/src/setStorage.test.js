import { config, DynamoDB } from "aws-sdk";
import uuid from "uuid";
import SetStorage from "./setStorage";

const tableName = "test";
let store;
let dynamo;

beforeEach(() => {
  config.update({ region: "eu-west-1" });

  const endpoint = "http://localhost:8000";
  dynamo = new DynamoDB({ endpoint: endpoint });

  store = new SetStorage({
    tableName: tableName,
    client: new DynamoDB.DocumentClient({ endpoint: endpoint })
  });

  return dynamo
    .createTable({
      TableName: tableName,
      AttributeDefinitions: [
        { AttributeName: "setNumber", AttributeType: "S" }
      ],
      KeySchema: [{ AttributeName: "setNumber", KeyType: "HASH" }],
      ProvisionedThroughput: {
        ReadCapacityUnits: 10,
        WriteCapacityUnits: 10
      }
    })
    .promise();
});

afterEach(() => {
  return dynamo.deleteTable({ TableName: tableName }).promise();
});

it("should write a single item", () => {
  const item = { setNumber: uuid(), value: 1 };

  return store
    .write(item)
    .then(() => store.read(item.setNumber))
    .then(item => expect(item).toEqual(item));
});

it("should handle a non-existing item", () =>
  store.read("123132").then(model => expect(model).toBeUndefined()));

it("should overwrite an existing item", () => {
  const item = { setNumber: uuid(), value: 1 };

  return Promise.resolve()
    .then(() => store.write(item))
    .then(() => store.write(item));
});
