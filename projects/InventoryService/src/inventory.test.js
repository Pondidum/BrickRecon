import Inventory from "./inventory";
import SetStorage from "./setStorage";

let api, setStorage, notifier, inventory;

beforeEach(() => {
  api = {
    getInventory: jest.fn()
  };
  setStorage = {
    read: jest.fn(),
    write: jest.fn()
  };
  setStorage.write.mockReturnValue(Promise.resolve());

  notifier = {
    publish: jest.fn()
  };
  notifier.publish.mockReturnValue(Promise.resolve());

  inventory = new Inventory(api, setStorage, notifier);
});

const modelInventory = [{ quantity: 5, boids: ["part"] }];

const storeDoesntContainTheModel = () =>
  setStorage.read.mockReturnValue(Promise.resolve());

const storeContainsTheModel = () =>
  setStorage.read.mockReturnValue(Promise.resolve({ setNumber: 1234 }));

const brickOwlDoesntHaveTheModel = () =>
  api.getInventory.mockReturnValue(Promise.resolve([]));

const brickOwlHasTheModel = () =>
  api.getInventory.mockReturnValue(Promise.resolve(modelInventory));
// {
//   owl.getBoid.mockReturnValue(Promise.resolve("abcde"));
//   owl.getModelInfo.mockReturnValue(Promise.resolve(model));
//   owl.getInventory.mockReturnValue(Promise.resolve(inv));
// };

// describe("buildModel", () => {
//   it("should contain the model info", () => {
//     brickOwlHasTheModel();

//     return inventory.buildModel(1234).then(result => {
//       expect(result).toMatchObject(model);
//       expect(result.inventory).toEqual(inv);
//     });
//   });

//   it("should return undefined for non existing model", () => {
//     brickOwlDoesntHaveTheModel();

//     return inventory
//       .buildModel(1234)
//       .then(model => expect(model).toBeUndefined());
//   });
// });

describe("updateInventory", () => {
  it("should not store anything if the set is unfound", () => {
    brickOwlDoesntHaveTheModel();
    storeDoesntContainTheModel();

    return inventory
      .updateInventory(1234)
      .then(() => expect(setStorage.write.mock.calls.length).toBe(0));
  });

  it("should do nothing if the model is in the store", () => {
    storeContainsTheModel();

    return inventory
      .updateInventory(1234)
      .then(() => expect(setStorage.write.mock.calls.length).toBe(0));
  });

  it("should store the model", () => {
    storeDoesntContainTheModel();
    brickOwlHasTheModel();

    return inventory.updateInventory(1234).then(() => {
      expect(setStorage.write.mock.calls.length).toBe(1);
      expect(notifier.publish.mock.calls.length).toBe(1);
    });
  });

  it("should reread the model if force is specified", () => {
    storeContainsTheModel();
    brickOwlHasTheModel();

    return inventory.updateInventory(1234, true).then(() => {
      expect(setStorage.write.mock.calls.length).toBe(1);
      expect(notifier.publish.mock.calls.length).toBe(1);
    });
  });
});
