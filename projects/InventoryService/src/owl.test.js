import Owl from "./owl";

let client, owl;

beforeEach(() => {
  client = jest.fn();
  owl = new Owl("TOKEN_WAT", client);
});

const mockResponse = res => client.mockReturnValue(Promise.resolve(res));
const expectUri = () => expect(client.mock.calls[0][0]);

describe("getBoid", () => {
  it("should lookup an id", () => {
    mockResponse({ boids: ["529600"] });

    return owl.getBoid(75042).then(id => {
      expectUri().toBe(
        "https://api.brickowl.com/v1/catalog/id_lookup?id=75042&key=TOKEN_WAT&type=Set"
      );
      expect(id).toEqual("529600");
    });
  });

  it("should return undefined for unrecognised number", () => {
    mockResponse({ boids: [] });

    return owl.getBoid(23131231).then(id => {
      expectUri().toBe(
        "https://api.brickowl.com/v1/catalog/id_lookup?id=23131231&key=TOKEN_WAT&type=Set"
      );
      expect(id).toBeUndefined();
    });
  });
});

describe("getInventory", () => {
  it("should fetch a set", () => {
    mockResponse({
      inventory: [
        { boid: "198888", quantity: "1", extra_quantity: "0", alt_link: 0 },
        { boid: "771344-81", quantity: "1", extra_quantity: "0", alt_link: 0 }
      ],
      unmatched_inventory: []
    });

    return owl.getInventory(98236).then(inventory => {
      expect(inventory).toEqual([
        { quantity: 1, boids: ["198888"] },
        { quantity: 1, boids: ["771344-81"] }
      ]);

      expectUri().toBe(
        "https://api.brickowl.com/v1/catalog/inventory?boid=98236&key=TOKEN_WAT"
      );
    });
  });

  it("should handle bad set ids", () => {
    mockResponse({ error: { status: "Invalid BOID" } });

    return owl.getInventory(12313131323).then(result => {
      expect(result).toBeUndefined();
      expectUri().toBe(
        "https://api.brickowl.com/v1/catalog/inventory?boid=12313131323&key=TOKEN_WAT"
      );
    });
  });

  it("should group alternalte parts", () => {
    mockResponse({
      inventory: [
        { boid: "nogroup-1", quantity: "1", alt_link: 0 },
        { boid: "group1-1", quantity: "1", alt_link: 1 },
        { boid: "group2-1", quantity: "1", alt_link: 2 },
        { boid: "group2-2", quantity: "1", alt_link: 2 },
        { boid: "nogroup-2", quantity: "1", alt_link: 0 },
        { boid: "group1-2", quantity: "1", alt_link: 1 }
      ]
    });

    return owl.getInventory(123).then(inventory => {
      expect(inventory).toEqual([
        { quantity: 1, boids: ["group1-1", "group1-2"] },
        { quantity: 1, boids: ["group2-1", "group2-2"] },
        { quantity: 1, boids: ["nogroup-1"] },
        { quantity: 1, boids: ["nogroup-2"] }
      ]);
      expectUri().toBe(
        "https://api.brickowl.com/v1/catalog/inventory?boid=123&key=TOKEN_WAT"
      );
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

    return owl.getModelInfo(529600).then(info => {
      expect(info.setNumbers).toEqual(["75042-1"]);
      expect(info.boid).toBe(529600);
      expect(info.name).toBe("LEGO Droid Gunship Set 75042");
      expect(info.url).toBe(
        "https://www.brickowl.com/catalog/lego-droid-gunship-set-75042"
      );

      expectUri().toBe(
        "https://api.brickowl.com/v1/catalog/lookup?boid=529600&key=TOKEN_WAT"
      );
    });
  });

  it("should handle bad set ids", () => {
    mockResponse({ error: { status: "Invalid BOID - 12313131323" } });

    return owl.getModelInfo(12313131323).then(result => {
      expect(result).toBeUndefined();
      expectUri().toBe(
        "https://api.brickowl.com/v1/catalog/lookup?boid=12313131323&key=TOKEN_WAT"
      );
    });
  });
});
