import Owl from "./owl";
import Storage from "./storage";
import Notifier from "./notifier";
import Inventory from "./inventory";

const tableName = process.env.TABLE_NAME;
const snsTopic = process.env.SNS_TOPIC;

const owl = new Owl();
const store = new Storage({ tableName: tableName, hashKey: "setNumber" });
const notifier = new Notifier(snsTopic);
const inventory = new Inventory(store, owl, notifier);

const handleSingle = record =>
  inventory.updateInventory(record.setNumber).catch(err => console.error(err));

exports.handler = (snsEvent, context, callback) => {
  const records = snsEvent.Records;

  const tasks = records
    .map(record => JSON.parse(record.Message))
    .filter(message => message.eventType === "MODEL_INVENTORY_REQUEST")
    .map(message => handleSingle(message));

  return Promise.all(records).then(() => callback());
};
