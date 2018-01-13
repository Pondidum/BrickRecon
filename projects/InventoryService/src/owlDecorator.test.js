import Decorator from "./owlDecorator";
import Owl from "./owl";

let client, decorator, cache;

beforeEach(() => {
  cache = { get: jest.fn(), write: jest.fn() };
  client = jest.fn();
  decorator = new Decorator(cache, new Owl("TOKEN_WAT", client));

  const cacheContent = {
    "324854-50": { color: 50, partNumber: "45" },
    "14440-38": { color: 38, partNumber: "46" },
    "954477-38": { color: 38, partNumber: "47" },
    "954477-12": { color: 12, partNumber: "47" },
    "529600": { color: 0, partNumber: "75042" }
  };
  cache.get.mockImplementation(key => Promise.resolve(cacheContent[key]));
  cache.write.mockReturnValue(Promise.resolve());
});

const mockResponse = res => client.mockReturnValue(Promise.resolve(res));

describe("getInventory", () => {
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
});

describe("getSetBoid", () => {
  it("should write the boid to the cache", () => {
    mockResponse({ boids: ["529600"] });

    return decorator.getSetBoid(75042).then(id => {
      expect(id).toEqual("529600");
      expect(cache.write.mock.calls[0]).toEqual([75042, "529600", 0]);
    });
  });

  it("should not write to the cache if unable to find the boid", () => {
    mockResponse({ boids: [] });

    return decorator.getSetBoid(75042).then(id => {
      expect(id).toBeUndefined();
      expect(cache.write.mock.calls.length).toBe(0);
    });
  });
});

describe("getModelInfo", () => {
  it("should return all set numbers", () => {
    mockResponse({
      boid: "529600",
      type: "Set",
      ids: [
        { id: "529600", type: "boid" },
        { id: "5702015119542", type: "ean" },
        { id: "6061124", type: "item_no" },
        { id: "673419210065", type: "upc" },
        { id: "75042-1", type: "set_number" },
        { id: "75042-1", type: "set_number" }
      ],
      name: "LEGO Droid Gunship Set 75042",
      url: "https://www.brickowl.com/catalog/lego-droid-gunship-set-75042",
      permalink: "https://www.brickowl.com/boid/529600",
      color_name: "",
      color_id: 0,
      color_hex: "000000",
      cat_name_path: "Sets / Star Wars / Episode III",
      missing_data: ""
    });

    return decorator.getModelInfo(529600).then(info => {
      expect(info.setNumber).toEqual("75042");
      expect(info).not.toHaveProperty("boid");
      expect(info.name).toBe("LEGO Droid Gunship Set 75042");
      expect(info.url).toBe(
        "https://www.brickowl.com/catalog/lego-droid-gunship-set-75042"
      );
    });
  });
});
