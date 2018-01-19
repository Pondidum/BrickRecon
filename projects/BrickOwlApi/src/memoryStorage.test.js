import MemoryStorage from "./memoryStorage";

let store;
beforeEach(() => (store = new MemoryStorage()));

it("should return empty object when there is nothing in the store", () =>
  store
    .getMany(["boid_1", "boid_2"])
    .then(result => expect(result).toEqual({})));

it("should write and read boids", () => {
  const input = {
    boid1: "part1",
    boid2: "part2",
    boid3: "part3"
  };

  return store
    .writeMany(input)
    .then(() =>
      store
        .getMany(Object.keys(input))
        .then(result => expect(result).toEqual(input))
    );
});

it("should only return queried boids", () => {
  const input = {
    boid1: "part1",
    boid2: "part2",
    boid3: "part3"
  };

  return store
    .writeMany(input)
    .then(() =>
      store
        .getMany(["boid2"])
        .then(result => expect(result).toEqual({ boid2: "part2" }))
    );
});

it("should not return non existing boids", () => {
  const input = {
    boid1: "part1",
    boid2: "part2",
    boid3: "part3"
  };

  return store
    .writeMany(input)
    .then(() =>
      store
        .getMany(["boid4"])
        .then(result => expect(result).not.toHaveProperty("boid4"))
    );
});
