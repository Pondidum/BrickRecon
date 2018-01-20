import MemoryStorage from "./memoryStorage";
import Api from "./api";

let storage, client, api;

beforeEach(() => {
  storage = new MemoryStorage();
  client = {
    boidFromSetNumber: () => Promise.resolve("set123"),
    getInventory: jest.fn(),
    getPartNumbers: jest.fn()
  };
  api = new Api({ storage: storage, client: client });

  client.getInventory.mockReturnValue(
    Promise.resolve([{ boids: ["543696-53"], quanity: 5 }])
  );
});

const expectInventoryToBeCorrect = inventory =>
  expect(inventory).toEqual([{ partNumber: "part1", quanity: 5, color: 53 }]);

it("should write to storage", () => {
  client.getPartNumbers.mockReturnValue(
    Promise.resolve({
      "543696-53": "part1"
    })
  );

  return api.getInventory("set123").then(inventory => {
    expectInventoryToBeCorrect(inventory);

    return storage
      .getMany(["543696-53"])
      .then(boids => expect(boids).toEqual({ "543696-53": "part1" }));
  });
});

it("should retrieve from storage", () => {
  client.getPartNumbers.mockReturnValue(Promise.resolve({}));

  return storage.writeMany({ "543696-53": "part1" }).then(() =>
    api.getInventory("set123").then(inventory => {
      expectInventoryToBeCorrect(inventory);
      expect(client.getPartNumbers.mock.calls.length).toEqual(0);
    })
  );
});
