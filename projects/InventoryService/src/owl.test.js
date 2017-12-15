import Owl from "./owl";

const owl = new Owl();

describe("getBoid", () => {
  it("should lookup an id", () =>
    owl.getBoid(75042).then(id => expect(id).toEqual("529600")));

  it("should return undefined for unrecognised number", () =>
    owl.getBoid(23131231).then(id => expect(id).toBeUndefined()));
});

describe("getInventory", () => {
  it("should fetch a set", () =>
    owl
      .getInventory(98236)
      .then(inventory =>
        expect(inventory).toEqual([
          { quantity: 1, boids: ["198888"] },
          { quantity: 1, boids: ["771344-81"] }
        ])
      ));

  it("should handle bad set ids", () =>
    owl
      .getInventory(12313131323)
      .then(result => expect(result).toBeUndefined()));

  it("should group alternalte parts", () => {
    const testFetch = () =>
      Promise.resolve({
        inventory: [
          { boid: "nogroup-1", quantity: "1", alt_link: 0 },
          { boid: "group1-1", quantity: "1", alt_link: 1 },
          { boid: "group2-1", quantity: "1", alt_link: 2 },
          { boid: "group2-2", quantity: "1", alt_link: 2 },
          { boid: "nogroup-2", quantity: "1", alt_link: 0 },
          { boid: "group1-2", quantity: "1", alt_link: 1 }
        ]
      });

    return new Owl(testFetch)
      .getInventory(123)
      .then(inventory =>
        expect(inventory).toEqual([
          { quantity: 1, boids: ["group1-1", "group1-2"] },
          { quantity: 1, boids: ["group2-1", "group2-2"] },
          { quantity: 1, boids: ["nogroup-1"] },
          { quantity: 1, boids: ["nogroup-2"] }
        ])
      );
  });
});
