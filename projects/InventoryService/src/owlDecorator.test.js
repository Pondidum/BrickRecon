import Decorator from "./owlDecorator";
import Owl from "./owl";

let client, decorator, cache;

beforeEach(() => {
  cache = { get: jest.fn() };
  client = jest.fn();
  decorator = new Decorator(cache, new Owl("TOKEN_WAT", client));

  const cacheContent = {
    "324854-50": { color: 50, partNumber: "45" },
    "14440-38": { color: 38, partNumber: "46" },
    "954477-38": { color: 38, partNumber: "47" },
    "954477-12": { color: 12, partNumber: "47" }
  };
  cache.get.mockImplementation(key => Promise.resolve(cacheContent[key]));
});

const mockResponse = res => client.mockReturnValue(Promise.resolve(res));

it("should replace brickowl ids", () => {
  mockResponse({
    inventory: [
      { boid: "324854-50", quantity: "4", alt_link: 0 },
      { boid: "14440-38", quantity: "1", alt_link: 0 },
      { boid: "954477-38", quantity: "4", alt_link: 0 },
      { boid: "954477-12", quantity: "2", alt_link: 0 }
    ],
    unmatched_inventory: []
  });

  return decorator
    .getInventory(98236)
    .then(inventory =>
      expect(inventory).toEqual([
        { quantity: 4, color: 50, partNumber: "45", otherPartNumbers: [] },
        { quantity: 1, color: 38, partNumber: "46", otherPartNumbers: [] },
        { quantity: 4, color: 38, partNumber: "47", otherPartNumbers: [] },
        { quantity: 2, color: 12, partNumber: "47", otherPartNumbers: [] }
      ])
    );
});

it("should replace within groups", () => {
  mockResponse({
    inventory: [
      { boid: "324854-50", quantity: "7", alt_link: 1 },
      { boid: "14440-38", quantity: "7", alt_link: 1 }
    ],
    unmatched_inventory: []
  });

  return decorator
    .getInventory(98236)
    .then(inventory =>
      expect(inventory).toEqual([
        { quantity: 7, color: 50, partNumber: "45", otherPartNumbers: ["46"] }
      ])
    );
});
