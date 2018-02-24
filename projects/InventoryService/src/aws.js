import { BrickOwlApi, DynamoStorage } from "brickowlapi";
import Notifier from "./notifier";
import Inventory from "./inventory";
import SetStorage from "./setStorage";

const brickOwlToken = process.env.BRICKOWL_TOKEN;
const setsTable = process.env.SETS_TABLE;
const boidsTable = process.env.BOIDS_TABLE;
const kafishLambda = process.env.KAFISH_LAMBDA;

const api = new BrickOwlApi({
  brickOwlToken: brickOwlToken,
  storage: new DynamoStorage(boidsTable)
});
const storage = new SetStorage({
  tableName: setsTable
});
const notifier = new Notifier({
  lambdaName: kafishLambda
});

const inventory = new Inventory(api, storage, notifier);

const isInventoryRequest = record =>
  record.MessageAttributes.EventType &&
  record.MessageAttributes.EventType.Value === "MODEL_INVENTORY_REQUEST";

const handleSingle = record =>
  inventory
    .updateInventory(record.setNumber, record.force)
    .catch(err => console.error(err));

exports.handler = (snsEvent, context, callback) => {
  const records = snsEvent.Records;

  const tasks = records
    .map(record => record.Sns)
    .filter(isInventoryRequest)
    .map(record => JSON.parse(record.Message))
    .map(message => handleSingle(message));

  return Promise.all(records).then(() => callback());
};
