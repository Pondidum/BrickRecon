import { config, DynamoDB } from "aws-sdk";
import uuid from "uuid";
import Storage from "./storage";

const tableName = "test";
let store;
let dynamo;

beforeEach(() => {
  config.update({ region: "eu-west-1" });

  const endpoint = "http://localhost:8000";

  store = new Storage({
    endpoint: endpoint,
    tableName: tableName,
    hashKey: "setNumber"
  });

  dynamo = new DynamoDB({
    endpoint: endpoint
  });

  return dynamo
    .createTable({
      TableName: tableName,
      AttributeDefinitions: [
        { AttributeName: "setNumber", AttributeType: "S" },
        { AttributeName: "value", AttributeType: "N" }
      ],
      KeySchema: [{ AttributeName: "setNumber", KeyType: "HASH" }],
      ProvisionedThroughput: {
        ReadCapacityUnits: 10,
        WriteCapacityUnits: 10
      },
      GlobalSecondaryIndexes: [
        {
          IndexName: "ByValue",
          KeySchema: [
            { AttributeName: "setNumber", KeyType: "HASH" },
            { AttributeName: "value", KeyType: "RANGE" }
          ],
          Projection: {
            ProjectionType: "ALL"
          },
          ProvisionedThroughput: {
            ReadCapacityUnits: 10,
            WriteCapacityUnits: 10
          }
        }
      ]
    })
    .promise();
});

afterEach(() => {
  return dynamo.deleteTable({ TableName: tableName }).promise();
});

it("should write multiple items", () => {
  const one = { setNumber: uuid(), value: 1 };
  const two = { setNumber: uuid(), value: 2 };

  return store
    .writeMany([one, two])
    .then(() =>
      store.read(one.setNumber).then(item => expect(item).toEqual(one))
    )
    .then(() =>
      store.read(two.setNumber).then(item => expect(item).toEqual(two))
    );
});

it("should write a single item", () => {
  const item = { setNumber: uuid(), value: 1 };

  return store
    .write(item)
    .then(() => store.read(item.setNumber))
    .then(item => expect(item).toEqual(item));
});

it("should read multiple items", () => {
  const item1 = { setNumber: uuid(), value: 1 };
  const item2 = { setNumber: uuid(), value: 3 };
  const item3 = { setNumber: uuid(), value: 7 };

  return store
    .writeMany([item1, item2, item3])
    .then(() =>
      store.readMany({
        expression: "#v < :v",
        nameMap: { "#v": "value" },
        valueMap: { ":v": 5 }
      })
    )
    .then(items => {
      items.sort((a, b) => a.value > b.value);
      expect(items).toEqual([item1, item2]);
    });
});
