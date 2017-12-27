import Inventory from "./inventory";
import Owl from "./owl";
let owl, store, notifier, inventory;

beforeEach(() => {
  owl = {
    getBoid: jest.fn(),
    getInventory: jest.fn(),
    getModelInfo: jest.fn()
  };
  store = {
    read: jest.fn(),
    write: jest.fn()
  };
  notifier = {
    publish: jest.fn()
  };

  store.write.mockReturnValue(Promise.resolve());
  notifier.publish.mockReturnValue(Promise.resolve());

  inventory = new Inventory(store, owl, notifier);
});

const model = {
  boid: "model_boid",
  name: "some model",
  url: "https://example.com",
  setNumbers: [123, 456]
};

const inv = [{ quantity: 5, boids: ["part"] }];

const storeDoesntContainTheModel = () =>
  store.read.mockReturnValue(Promise.resolve());

const storeContainsTheModel = () =>
  store.read.mockReturnValue(Promise.resolve({ setNumber: 1234 }));

const brickOwlDoesntHaveTheModel = () =>
  owl.getBoid.mockReturnValue(Promise.resolve());

const brickOwlHasTheModel = () => {
  owl.getBoid.mockReturnValue(Promise.resolve("abcde"));
  owl.getModelInfo.mockReturnValue(Promise.resolve(model));
  owl.getInventory.mockReturnValue(Promise.resolve(inv));
};

describe("buildModel", () => {
  it("should contain the model info", () => {
    brickOwlHasTheModel();

    return inventory.buildModel(1234).then(result => {
      expect(result).toMatchObject(model);
      expect(result.inventory).toEqual(inv);
    });
  });

  it("should return undefined for non existing model", () => {
    brickOwlDoesntHaveTheModel();

    return inventory
      .buildModel(1234)
      .then(model => expect(model).toBeUndefined());
  });
});

describe("updateInventory", () => {
  it("should not store anything if the set is unfound", () => {
    brickOwlDoesntHaveTheModel();
    storeDoesntContainTheModel();

    return inventory
      .updateInventory(1234)
      .then(() => expect(store.write.mock.calls.length).toBe(0));
  });

  it("should do nothing if the model is in the store", () => {
    storeContainsTheModel();

    return inventory
      .updateInventory(1234)
      .then(() => expect(store.write.mock.calls.length).toBe(0));
  });

  it("should store the model", () => {
    storeDoesntContainTheModel();
    brickOwlHasTheModel();

    return inventory.updateInventory(1234).then(() => {
      expect(store.write.mock.calls.length).toBe(1);
      expect(notifier.publish.mock.calls.length).toBe(1);
    });
  });

  it("should reread the model if force is specified", () => {
    storeContainsTheModel();
    brickOwlHasTheModel();

    return inventory.updateInventory(1234, true).then(() => {
      expect(store.write.mock.calls.length).toBe(1);
      expect(notifier.publish.mock.calls.length).toBe(1);
    });
  });
});
