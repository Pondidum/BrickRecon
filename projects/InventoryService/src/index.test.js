import { getBoid, getInventory } from "./index";

describe("getBoid", () => {
  it("should lookup an id", () => {
    return getBoid(75042).then(id => expect(id).toEqual("529600"));
  });

  it("should return undefined for unrecognised number", () => {
    return getBoid(23131231).then(id => expect(id).toBeUndefined());
  });
});

describe("getInventory", () => {
  it("should fetch a set", () => {
    return getInventory(98236).then(inventory => {
      expect(inventory).toEqual([
        { quantity: 1, boids: ["198888"] },
        { quantity: 1, boids: ["771344-81"] }
      ]);
    });
  });

  it("should handle bad set ids", () => {
    return getInventory(12313131323).then(result => {
      expect(result).toBeUndefined();
    });
  });
});
