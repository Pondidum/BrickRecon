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

exports.handler = (event, context, callback) =>
  inventory
    .updateInventory(event.setNumber)
    .then(() => callback())
    .catch(err => {
      console.error(err);
      callback(err, err.toString());
    });
