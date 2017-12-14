import { getBoid } from "./index";

it("should lookup an id", () => {
  return getBoid(75042).then(id => expect(id).toEqual(529600));
});

it("should return undefined for unrecognised number", () => {
  return getBoid(23131231).then(id => expect(id).toBeUndefined());
});
