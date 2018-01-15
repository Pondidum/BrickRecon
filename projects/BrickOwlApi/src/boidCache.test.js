import BoidCache from "./boidCache";

let client, storage, cache;

beforeEach(() => {
  storage = {
    getMany: jest.fn(),
    writeMany: jest.fn()
  };

  client = {
    getPartNumbers: jest.fn()
  };

  storage.writeMany.mockReturnValue(Promise.resolve());

  cache = new BoidCache(storage, client);
});

const storageReturns = items =>
  storage.getMany.mockReturnValue(Promise.resolve(items));
const clientReturns = items =>
  client.getPartNumbers.mockReturnValue(Promise.resolve(items));

it("should items missing from storage to storage", () => {
  storageReturns({});
  clientReturns({ boid1: "part1", boid2: "part2" });

  return cache.getMany(["boid1", "boid2"]).then(map => {
    expect(map).toEqual({ boid1: "part1", boid2: "part2" });
    expect(storage.writeMany.mock.calls.length).toEqual(1);
    expect(storage.writeMany.mock.calls[0][0]).toEqual({
      boid1: "part1",
      boid2: "part2"
    });
  });
});

it("should return from storage if storage has all items", () => {
  storageReturns({ boid1: "part1", boid2: "part2" });
  clientReturns({ boid1: "client1", boid2: "client2" });

  return cache.getMany(["boid1", "boid2"]).then(map => {
    expect(map).toEqual({ boid1: "part1", boid2: "part2" });
  });
});

it("should return from storage and client when storage doesnt have everything", () => {
  storageReturns({ boid1: "part1" });
  clientReturns({ boid2: "client2" });

  return cache.getMany(["boid1", "boid2"]).then(map => {
    expect(map).toEqual({ boid1: "part1", boid2: "client2" });
    expect(storage.writeMany.mock.calls[0][0]).toEqual({ boid2: "client2" });
  });
});
