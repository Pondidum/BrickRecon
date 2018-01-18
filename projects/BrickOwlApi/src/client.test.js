import Client from "./client";
import { mapFrom } from "./util";

let fetcher, client;

beforeEach(() => {
  fetcher = jest.fn();
  client = new Client("TOKEN_WAT", {
    fetcher: fetcher,
    batchSize: 5
  });
});

const mockResponse = res => fetcher.mockReturnValue(Promise.resolve(res));

describe("getSetBoid", () => {
  it("should read the response", () => {
    mockResponse({ boids: ["529600"] });
    return client.getSetBoid(75042).then(id => expect(id).toBe("529600"));
  });

  it("should an handle unrecognised boid", () => {
    mockResponse({ boids: [] });
    return client.getSetBoid(75042).then(id => expect(id).toBe(undefined));
  });
});

describe("getInventory", () => {
  it("should handle bad set ids", () => {
    mockResponse({ error: { status: "Invalid BOID" } });

    return client
      .getInventory("529600")
      .then(inventory => expect(inventory).toEqual([]));
  });

  it("should return the inventory", () => {
    mockResponse({
      inventory: [
        { boid: "324854-50", quantity: "4", extra_quantity: "0", alt_link: 0 },
        { boid: "14440-38", quantity: "1", extra_quantity: "0", alt_link: 0 },
        { boid: "324854-38", quantity: "2", extra_quantity: "0", alt_link: 0 }
      ]
    });

    return client
      .getInventory("529600")
      .then(inventory =>
        expect(inventory).toEqual([
          { boids: ["324854-50"], quantity: 4 },
          { boids: ["14440-38"], quantity: 1 },
          { boids: ["324854-38"], quantity: 2 }
        ])
      );
  });

  it("should group alternate parts", () => {
    mockResponse({
      inventory: [
        { boid: "11111-50", quantity: "4", extra_quantity: "0", alt_link: 1 },
        { boid: "22222-38", quantity: "1", extra_quantity: "0", alt_link: 0 },
        { boid: "33333-38", quantity: "4", extra_quantity: "0", alt_link: 1 }
      ]
    });

    return client
      .getInventory("529600")
      .then(inventory =>
        expect(inventory).toEqual([
          { boids: ["11111-50", "33333-38"], quantity: 4 },
          { boids: ["22222-38"], quantity: 1 }
        ])
      );
  });
});

describe("getPartNumbers", () => {
  const defaultResponse = {
    items: {
      "543696-53": {
        boid: "543696-53",
        type: "Part",
        ids: [
          { id: "30136", type: "design_id" },
          { id: "30136", type: "design_id" },
          { id: "30136", type: "ldraw" },
          { id: "30136", type: "peeron_id" },
          { id: "4114054", type: "item_no" },
          { id: "4114054", type: "item_no" },
          { id: "543696-53", type: "boid" }
        ]
      },
      "297160-53": {
        boid: "297160-53",
        type: "Part",
        ids: [
          { id: "297160-53", type: "boid" },
          { id: "30132", type: "design_id" },
          { id: "30132", type: "design_id" },
          { id: "30132", type: "ldraw" },
          { id: "30132", type: "peeron_id" },
          { id: "4105286", type: "item_no" },
          { id: "4105286", type: "item_no" }
        ]
      }
    }
  };

  it("should do nothing when empty array passed", () => {
    mockResponse(defaultResponse);
    return client.getPartNumbers([]).then(parts => expect(parts).toEqual({}));
  });

  it("should send one request for less than the chunk size", () => {
    mockResponse(defaultResponse);
    return client.getPartNumbers(["543696-53", "297160-53"]).then(parts => {
      expect(parts).toEqual({
        "543696-53": "30136",
        "297160-53": "30132"
      });

      expect(fetcher.mock.calls.length).toBe(1);
    });
  });

  it("should send one request for less than the chunk size", () => {
    const boids = [...new Array(13).keys()].map(id => id.toString());

    const items = mapFrom(
      boids,
      id => id,
      id => ({ boid: id, ids: [{ id: "part:" + id, type: "ldraw" }] })
    );

    mockResponse({ items: items });

    return client.getPartNumbers(boids).then(parts => {
      expect(Object.keys(parts).length).toBe(boids.length);
      expect(fetcher.mock.calls.length).toBe(3);
    });
  });
});
