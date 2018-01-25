import { BrickOwlApi, MemoryStorage } from "brickowlapi";
import Inventory from "./inventory";

const brickOwlToken = process.env.BRICKOWL_TOKEN;

it("should really work", () => {
  const boidStorage = new MemoryStorage();
  const api = new BrickOwlApi({
    brickOwlToken: brickOwlToken,
    storage: boidStorage
  });
  const storage = {
    read: () => Promise.resolve(),
    write: items => {
      console.log("write", items);
      return Promise.resolve();
    }
  };
  const notifier = {
    publish: () => Promise.resolve()
  };

  const inventory = new Inventory(api, storage, notifier);

  return inventory.updateInventory("75042", false);
});
