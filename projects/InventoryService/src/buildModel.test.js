import buildModel from "./buildModel";

const model = {
  boid: "model_boid",
  name: "some model",
  url: "https://example.com",
  setNumbers: [123, 456]
};

const inv = [{ quantity: 5, boids: ["part"] }];

const owl = {
  getBoid: setNumber => Promise.resolve("abcde"),
  getModelInfo: boid => Promise.resolve(model),
  getInventory: boid => Promise.resolve(inv)
};

it("should contain the model info", () =>
  buildModel(owl, 1234).then(model => expect(model).toMatchObject(model)));

it("should contain the inventory", () =>
  buildModel(owl, 1234).then(model => expect(model.inventory).toEqual(inv)));
