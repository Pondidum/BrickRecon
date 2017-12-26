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

  inventory = new Inventory(store, owl, notifier);
});

describe("buildModel", () => {
  const model = {
    boid: "model_boid",
    name: "some model",
    url: "https://example.com",
    setNumbers: [123, 456]
  };

  const inv = [{ quantity: 5, boids: ["part"] }];

  it("should contain the model info", () => {
    owl.getBoid.mockReturnValue(Promise.resolve("abcde"));
    owl.getModelInfo.mockReturnValue(Promise.resolve(model));
    owl.getInventory.mockReturnValue(Promise.resolve(inv));

    return inventory.buildModel(1234).then(result => {
      expect(result).toMatchObject(model);
      expect(result.inventory).toEqual(inv);
    });
  });

  it("should return undefined for non existing model", () => {
    owl.getBoid.mockReturnValue(Promise.resolve());

    return inventory
      .buildModel(1234)
      .then(model => expect(model).toBeUndefined());
  });
});

describe("updateInventory", () => {
  it("should not store anything if the set is unfound", () => {
    owl.getBoid.mockReturnValue(Promise.resolve());
    store.read.mockReturnValue(Promise.resolve());

    return inventory
      .updateInventory(1234)
      .then(() => expect(store.write.mock.calls.length).toBe(0));
  });

  it("should do nothing if the model is in the store", () => {
    store.read.mockReturnValue(Promise.resolve({ setNumber: 1234 }));

    return inventory
      .updateInventory(1234)
      .then(() => expect(store.write.mock.calls.length).toBe(0));
  });

  it("should store the model", () => {
    store.read.mockReturnValue(Promise.resolve());
    store.write.mockReturnValue(Promise.resolve());

    owl.getBoid.mockReturnValue(Promise.resolve("boid"));
    owl.getModelInfo.mockReturnValue(Promise.resolve({ name: "some set" }));
    owl.getInventory.mockReturnValue(Promise.resolve([]));

    notifier.publish.mockReturnValue(Promise.resolve());

    return inventory.updateInventory(1234).then(() => {
      expect(store.write.mock.calls.length).toBe(1);
      expect(notifier.publish.mock.calls.length).toBe(1);
    });
  });
});
