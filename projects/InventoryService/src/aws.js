import { BrickOwlApi, DynamoStorage } from "brickowlapi";
import Notifier from "./notifier";
import Inventory from "./inventory";
import SetStorage from "./setStorage";

const brickowlToken = process.env.BRICKOWL_TOKEN;
const setsTable = process.env.SETS_TABLE;
const boidsTable = process.env.BOIDS_TABLE;
const snsTopic = process.env.SNS_TOPIC;

const api = new BrickOwlApi({
  brickowlToken: brickowlToken,
  storage: new DynamoStorage(boidsTable)
});
const storage = new SetStorage({
  tableName: setsTable
});
const notifier = new Notifier(snsTopic);

const inventory = new Inventory(api, storage, notifier);

const handleSingle = record =>
  inventory
    .updateInventory(record.setNumber, record.force)
    .catch(err => console.error(err));

exports.handler = (snsEvent, context, callback) => {
  const records = snsEvent.Records;

  const tasks = records
    .map(record => JSON.parse(record.Sns.Message))
    .filter(message => message.eventType === "MODEL_INVENTORY_REQUEST")
    .map(message => handleSingle(message));

  return Promise.all(records).then(() => callback());
};
