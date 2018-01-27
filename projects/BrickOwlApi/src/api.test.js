import MemoryStorage from "./memoryStorage";
import DynamoStorage from "./dynamoStorage";
import Api from "./api";
import AWS, { DynamoDB } from "aws-sdk";

let storage, client, api;

beforeEach(() => {
  storage = new MemoryStorage();
  client = {
    getSetBoid: () => Promise.resolve("set123"),
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
      "543696": "part1"
    })
  );

  return api.getInventory("set123").then(inventory => {
    expectInventoryToBeCorrect(inventory);

    return storage
      .getMany(["543696"])
      .then(boids => expect(boids).toEqual({ "543696": "part1" }));
  });
});

it("should retrieve from storage", () => {
  client.getPartNumbers.mockReturnValue(Promise.resolve({}));

  return storage.writeMany({ "543696": "part1" }).then(() =>
    api.getInventory("set123").then(inventory => {
      expectInventoryToBeCorrect(inventory);
      expect(client.getPartNumbers.mock.calls.length).toEqual(0);
    })
  );
});

it("should handle the set not existing", () => {
  client.getSetBoid = () => Promise.resolve();

  return api
    .getInventory("set123")
    .then(inventory => expect(inventory).toEqual([]));
});

it("should work for real", () => {
  AWS.config.update({ region: "eu-west-1" });

  const real = new Api({
    brickOwlToken: process.env.BRICKOWL_TOKEN,
    storage: new DynamoStorage("boidlookup_test", {
      client: new DynamoDB.DocumentClient({
        endpoint: "http://localhost:8000"
      })
    })
  });

  return real.getInventory("75042").then(inv => {
    var badParts = inv.filter(part => part.partNumber === undefined);
    expect(badParts.length).toEqual(0);
  });
});
