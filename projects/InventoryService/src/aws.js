import Owl from "./owl";
import Storage from "./storage";
import Notifier from "./notifier";
import Inventory from "./inventory";

const brickowlToken = process.env.BRICKOWL_TOKEN;
const tableName = process.env.TABLE_NAME;
const snsTopic = process.env.SNS_TOPIC;

const owl = new Owl(brickowlToken);
const store = new Storage({ tableName: tableName, hashKey: "setNumber" });
const notifier = new Notifier(snsTopic);
const inventory = new Inventory(store, owl, notifier);

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
