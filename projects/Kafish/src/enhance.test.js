import enhance from "./enhance";

it("should add eventId and timestamp if not present", () => {
  const event = enhance({
    wat: "r u doin"
  });

  expect(event.eventId).toBeDefined();
  expect(event.timestamp).toBeDefined();
});

it("should not replace an existing eventId", () => {
  const event = enhance({
    eventId: "existing",
    wat: "r u doin"
  });

  expect(event.eventId).toEqual("existing");
});

it("should replace an existing eventId", () => {
  const event = enhance({
    wat: "r u doin",
    timestamp: "1234"
  });

  expect(event.timestamp).not.toEqual("1234");
});

it("should replace empty string properties with null", () => {
  const event = enhance({
    wat: "",
    again: "nope"
  });

  expect(event.wat).toEqual(null);
  expect(event.again).toEqual("nope");
});
