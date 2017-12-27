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
        { boid: "324854-50", quantity: "4", extra_quantity: "0", alt_link: 0 },
        { boid: "14440-38", quantity: "1", extra_quantity: "0", alt_link: 0 },
        { boid: "954477-38", quantity: "4", extra_quantity: "0", alt_link: 0 },
        { boid: "739988-97", quantity: "1", extra_quantity: "0", alt_link: 0 },
        { boid: "267724-38", quantity: "2", extra_quantity: "0", alt_link: 0 },
        { boid: "408757-38", quantity: "2", extra_quantity: "0", alt_link: 0 },
        { boid: "369885-38", quantity: "8", extra_quantity: "0", alt_link: 0 },
        { boid: "487300-38", quantity: "4", extra_quantity: "0", alt_link: 0 },
        { boid: "44432-38", quantity: "2", extra_quantity: "0", alt_link: 0 },
        { boid: "886867-38", quantity: "2", extra_quantity: "0", alt_link: 0 }
      ],
      unmatched_inventory: []
    });

    return owl.getInventory(98236).then(inventory => {
      expect(inventory).toEqual([
        { quantity: 4, boids: ["324854-50"] },
        { quantity: 1, boids: ["14440-38"] },
        { quantity: 4, boids: ["954477-38"] },
        { quantity: 1, boids: ["739988-97"] },
        { quantity: 2, boids: ["267724-38"] },
        { quantity: 2, boids: ["408757-38"] },
        { quantity: 8, boids: ["369885-38"] },
        { quantity: 4, boids: ["487300-38"] },
        { quantity: 2, boids: ["44432-38"] },
        { quantity: 2, boids: ["886867-38"] }
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
        { boid: "ng-1", quantity: "1", alt_link: 0 },
        { boid: "g1-1", quantity: "1", alt_link: 1 },
        { boid: "g2-1", quantity: "1", alt_link: 2 },
        { boid: "g2-2", quantity: "1", alt_link: 2 },
        { boid: "ng-2", quantity: "1", alt_link: 0 },
        { boid: "g1-2", quantity: "1", alt_link: 1 }
      ]
    });

    return owl.getInventory(123).then(inventory => {
      expect(inventory).toEqual([
        { quantity: 1, boids: ["g1-1", "g1-2"] },
        { quantity: 1, boids: ["g2-1", "g2-2"] },
        { quantity: 1, boids: ["ng-1"] },
        { quantity: 1, boids: ["ng-2"] }
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
