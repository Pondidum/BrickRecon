import rebuild from "./rebuild";

const image = {
  Message: {
    S: "New item!"
  },
  Id: {
    N: "101"
  }
};

it("should be able to reconstruct a record from dynamo db", () => {
  expect(rebuild(image)).toEqual({
    Message: "New item!",
    Id: 101
  });
});
