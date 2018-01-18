import { mapFrom } from "./util";

describe("mapFrom", () => {
  it("should return empty object for empty array", () =>
    expect(mapFrom([])).toEqual({}));

  it("should return the item as the value if no valueFunc specified", () => {
    const input = [{ id: "wat" }, { id: "is" }, { id: "this" }];
    const map = mapFrom(input, x => x.id);

    expect(map).toEqual({
      wat: { id: "wat" },
      is: { id: "is" },
      this: { id: "this" }
    });
  });

  it("should build a map", () => {
    const input = [
      { id: "a", name: "one" },
      { id: "b", name: "two" },
      { id: "c", name: "three" }
    ];

    const map = mapFrom(input, x => x.id, x => x.name);

    expect(map).toEqual({
      a: "one",
      b: "two",
      c: "three"
    });
  });
});
